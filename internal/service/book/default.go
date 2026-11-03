package book

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/api/protocol"
	"github.com/tp86/legimi-go/internal/model"
	"github.com/tp86/legimi-go/internal/service"
)

type defaultBookService struct {
	sessionService    service.Session
	client            api.Client
	downloadPresenter service.DownloadPresenter
}

func (bs defaultBookService) ListBooks() ([]model.BookMetadata, error) {
	// TODO better error handling
	sessionId, err := bs.sessionService.GetSession()
	if err != nil {
		return nil, err
	}
	list := make([]model.BookMetadata, 0)
	request := model.NewBookListRequest(sessionId)
	var bookList model.BookList
	for {
		err := bs.client.Exchange(request, &bookList)
		if err != nil {
			return list, err
		}
		if len(bookList) == 0 {
			break
		}
		for _, book := range bookList {
			list = append(list, book)
		}
		request.NextPage = bookList[len(bookList)-1].NextPage
	}
	return list, nil
}

func (bs defaultBookService) DownloadBooks(bookIds []uint64) error {
	// TODO concurrent download of all books
	errs := make([]error, 0)
	for _, id := range bookIds {
		errs = append(errs, bs.downloadBook(id))
	}
	return errors.Join(errs...)
}

func (bs defaultBookService) downloadBook(id uint64) error {
	// TODO concurrent downloader
	sessionId, err := bs.sessionService.GetSession()
	if err != nil {
		return err
	}
	book, err := bs.getBookMetadata(sessionId, id)
	if err != nil {
		return err
	}
	bookDownloadDetails, err := bs.getBookDownloadDetails(sessionId, book)
	if err != nil {
		return err
	}
	return bs.download(book, bookDownloadDetails)
}

func (bs defaultBookService) getBookMetadata(sessionId string, bookId uint64) (model.BookMetadata, error) {
	metadataRequest := model.NewBookListRequest(sessionId)
	metadataRequest.BookId = bookId
	var bookList model.BookList
	err := bs.client.Exchange(metadataRequest, &bookList)
	if err != nil {
		return model.BookMetadata{}, err
	}
	if len(bookList) != 1 {
		return model.BookMetadata{}, fmt.Errorf("unexpected book metadata list count: %d, expected 1", len(bookList))
	}
	return bookList[0], nil
}

const maxDownloadDetailsGetAttempts = 5

func (bs defaultBookService) getBookDownloadDetails(sessionId string, book model.BookMetadata) (model.BookDownloadDetails, error) {
	downloadDetailsRequest := model.NewBookDownloadDetailsRequest(sessionId, book.Id, book.Version)
	var bookDownloadDetails model.BookDownloadDetails
	// TODO refactor & test
	attempt := 0
	for ; attempt < maxDownloadDetailsGetAttempts; attempt++ {
		if err := bs.client.Exchange(downloadDetailsRequest, &bookDownloadDetails); err != nil {
			if err, ok := err.(protocol.ErrorResponse); ok && err.Type == protocol.BookDownloadDetailsPreparingError {
				// special case - download details are being prepared, try to repeat after some time
				bs.downloadPresenter.Wait(book)
				time.Sleep(2 * time.Second)
				continue
			}
			return bookDownloadDetails, err
		}
		break
	}
	if attempt == maxDownloadDetailsGetAttempts {
		return bookDownloadDetails, fmt.Errorf("couldn't get download details after %d attempts, try downloading book again after some time", attempt)
	}
	return bookDownloadDetails, nil
}

const (
	downloadChunkSize                  uint64 = 81920
	maxDownloadAttemptsWithoutProgress        = 5
)

// variable, so that tests don't have to wait
var downloadRetryDelay = 2 * time.Second

func (bs defaultBookService) download(book model.BookMetadata, downloadDetails model.BookDownloadDetails) error {
	bs.downloadPresenter.Start(book)
	fileName := fmt.Sprintf("%d.mobi", book.Id)
	// download to temporary file, so that incomplete download never looks like a book
	partFileName := fileName + ".part"
	file, err := os.Create(partFileName)
	if err != nil {
		return err
	}
	err = bs.downloadTo(file, book, downloadDetails)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Rename(partFileName, fileName)
	}
	if err != nil {
		os.Remove(partFileName)
		return err
	}
	bs.downloadPresenter.End(book)
	return nil
}

func (bs defaultBookService) downloadTo(file *os.File, book model.BookMetadata, downloadDetails model.BookDownloadDetails) error {
	client := http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest(http.MethodGet, downloadDetails.Url, nil)
	if err != nil {
		return err
	}
	var downloadedBytes uint64 = 0
	attemptsWithoutProgress := 0
	for downloadedBytes < downloadDetails.Size {
		if attemptsWithoutProgress > 0 {
			time.Sleep(downloadRetryDelay)
		}
		// always continue from bytes actually received, chunk may come back short
		request.Header.Set("range", fmt.Sprintf("bytes=%d-%d", downloadedBytes, min(downloadedBytes+downloadChunkSize, downloadDetails.Size)-1))
		bytesRead, err := downloadChunk(client, request, file, downloadedBytes, downloadDetails.Size)
		downloadedBytes += uint64(bytesRead)
		if bytesRead > 0 {
			attemptsWithoutProgress = 0
			bs.downloadPresenter.Part(book)
			continue
		}
		if err, ok := err.(fatalDownloadError); ok {
			return err
		}
		attemptsWithoutProgress++
		if attemptsWithoutProgress == maxDownloadAttemptsWithoutProgress {
			if err == nil {
				err = errors.New("no data received")
			}
			return fmt.Errorf("download of book %d stalled at %d of %d bytes: %w", book.Id, downloadedBytes, downloadDetails.Size, err)
		}
	}
	if downloadedBytes != downloadDetails.Size {
		return fmt.Errorf("downloaded %d bytes of book %d, expected %d", downloadedBytes, book.Id, downloadDetails.Size)
	}
	return nil
}

// fatalDownloadError is returned for errors that wouldn't go away by repeating request
type fatalDownloadError struct {
	error
}

// downloadChunk writes response to requested range into file and returns number of bytes written.
// Returned bytes count may be non-zero along with error if response was cut off, bytes should be kept then.
func downloadChunk(client http.Client, request *http.Request, file *os.File, offset, size uint64) (int64, error) {
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	switch {
	case response.StatusCode == http.StatusPartialContent:
		if err := checkContentRange(response.Header.Get("content-range"), offset, size); err != nil {
			return 0, fatalDownloadError{err}
		}
	case response.StatusCode == http.StatusOK && offset == 0:
		// server ignored range and sends whole file
	case response.StatusCode >= 500:
		return 0, fmt.Errorf("unexpected response status: %s", response.Status)
	default:
		return 0, fatalDownloadError{fmt.Errorf("unexpected response status: %s", response.Status)}
	}
	return file.ReadFrom(response.Body)
}

func checkContentRange(contentRange string, offset, size uint64) error {
	if contentRange == "" {
		return nil
	}
	var start, end uint64
	var total string
	if _, err := fmt.Sscanf(contentRange, "bytes %d-%d/%s", &start, &end, &total); err != nil {
		return fmt.Errorf("invalid content-range in response: %q", contentRange)
	}
	if start != offset {
		return fmt.Errorf("content-range in response starts at %d, requested %d", start, offset)
	}
	if total != "*" && total != fmt.Sprint(size) {
		return fmt.Errorf("content-range in response has total size %s, expected %d", total, size)
	}
	return nil
}
