package presenter

import (
	"fmt"

	"github.com/tp86/legimi-go/internal/model"
)

type defaultBookListPresenter struct{}

func (defaultBookListPresenter) Present(bookList []model.BookMetadata) {
	for _, book := range bookList {
		fmt.Printf("%8d: \"%s\", %s, downloaded: %t\n", book.Id, book.Title, book.Author, book.Downloaded)
	}
}

type defaultBookDownloadPresenter struct{}

func (defaultBookDownloadPresenter) Start(book model.BookMetadata) {
	fmt.Printf("Downloading book %d: \"%s\" ", book.Id, book.Title)
}

func (defaultBookDownloadPresenter) Part(book model.BookMetadata) {
	fmt.Print(".")
}

func (defaultBookDownloadPresenter) End(book model.BookMetadata) {
	fmt.Println(" done")
}
