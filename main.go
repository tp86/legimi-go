package main

import (
	"fmt"

	"github.com/tp86/legimi-go/internal/service"
	"github.com/tp86/legimi-go/internal/usecase"
)

var (
	bookLister        usecase.BookLister
	bookListPresenter service.BookListPresenter
)

func main() {
	bookList, err := bookLister.ListBooks()
	if err != nil {
		fmt.Println(err)
	}
	bookListPresenter.Present(bookList)
}
