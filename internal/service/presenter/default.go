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
