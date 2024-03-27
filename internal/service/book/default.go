package book

import (
	"errors"
	"fmt"

	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/model"
	"github.com/tp86/legimi-go/internal/service"
)

type defaultBookService struct {
	sessionService service.Session
	client         api.Client
}

func (bs defaultBookService) ListBooks() error {
	// TODO better error handling
	sessionId, err := bs.sessionService.GetSession()
	if err != nil {
		return err
	}
	list := make([]string, 0)
	request := model.NewBookListRequest(sessionId)
	var bookList model.BookList
	for {
		err := bs.client.Exchange(request, &bookList)
		if err != nil {
			return err
		}
		if len(bookList) == 0 {
			break
		}
		for _, book := range bookList {
			list = append(list, book.String())
		}
		request.NextPage = bookList[len(bookList)-1].NextPage
	}
	for _, entry := range list {
		fmt.Println(entry)
	}
	return nil
}

func (bs defaultBookService) DownloadBooks(bookIds []uint64) error {
	// TODO parallel download
	errs := make([]error, 0)
	for _, id := range bookIds {
		errs = append(errs, bs.downloadBook(id))
	}
	return errors.Join(errs...)
}

func (bs defaultBookService) downloadBook(id uint64) error {
	return nil
}
