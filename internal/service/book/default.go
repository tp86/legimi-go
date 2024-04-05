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

const downloadChunkSize uint64 = 81920

func (bs defaultBookService) download(book model.BookMetadata, downloadDetails model.BookDownloadDetails) error {
	bs.downloadPresenter.Start(book)
	file, err := os.Create(fmt.Sprintf("%d.mobi", book.Id))
	if err != nil {
		return err
	}
	defer file.Close()
	client := http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest(http.MethodGet, downloadDetails.Url, nil)
	if err != nil {
		return err
	}
	var downloadedBytes uint64 = 0
	for i := uint64(0); downloadedBytes < downloadDetails.Size; i++ {
		request.Header.Set("range", fmt.Sprintf("bytes=%d-%d", i*downloadChunkSize, min((i+1)*downloadChunkSize-1, downloadDetails.Size)))
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		bytesRead, err := file.ReadFrom(response.Body)
		response.Body.Close()
		if err != nil {
			return err
		}
		downloadedBytes += uint64(bytesRead)
		bs.downloadPresenter.Part(book)
	}
	bs.downloadPresenter.End(book)
	return nil
}
