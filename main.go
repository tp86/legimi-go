package main

import (
	"fmt"

	"github.com/tp86/legimi-go/internal/service"
	"github.com/tp86/legimi-go/internal/usecase"
)

var (
	bookLister        usecase.BookLister
	bookListPresenter service.BookListPresenter
	bookDownloader    usecase.BookDownloader
)

func main() {
	bookList, err := bookLister.ListBooks()
	if err != nil {
		fmt.Println(err)
	}
	bookListPresenter.Present(bookList)
	bookId := 1022214
	err = bookDownloader.DownloadBooks([]uint64{uint64(bookId)})
	if err != nil {
		fmt.Println(err)
	}
}
