package book

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/model"
	"github.com/tp86/legimi-go/internal/service"
)

type defaultBookService struct {
	sessionService service.Session
	client         api.Client
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

const downloadChunkSize uint64 = 81920

func (bs defaultBookService) downloadBook(id uint64) error {
	// TODO refactor
	// TODO concurrent downloader
	sessionId, err := bs.sessionService.GetSession()
	if err != nil {
		return err
	}
	metadataRequest := model.NewBookListRequest(sessionId)
	metadataRequest.BookId = id
	var bookList model.BookList
	err = bs.client.Exchange(metadataRequest, &bookList)
	if err != nil {
		return err
	}
	if len(bookList) != 1 {
		return fmt.Errorf("unexpected book metadata list count: %d, expected 1", len(bookList))
	}
	book := bookList[0]
	downloadDetailsRequest := model.NewBookDownloadDetailsRequest(sessionId, book.Id, book.Version)
	var bookDownloadDetails model.BookDownloadDetails
	err = bs.client.Exchange(downloadDetailsRequest, &bookDownloadDetails)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodGet, bookDownloadDetails.Url, nil)
	if err != nil {
		return err
	}
	client := http.Client{Timeout: 30 * time.Second}
	file, err := os.Create(fmt.Sprintf("%d.mobi", book.Id))
	if err != nil {
		return err
	}
	defer file.Close()
	var downloadedBytes uint64 = 0
	for i := uint64(0); downloadedBytes < bookDownloadDetails.Size; i++ {
		request.Header.Set("range", fmt.Sprintf("bytes=%d-%d", i*downloadChunkSize, min((i+1)*downloadChunkSize-1, bookDownloadDetails.Size)))
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
	}
	return nil
}
