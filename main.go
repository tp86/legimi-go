package main

import (
	"fmt"

	"github.com/tp86/legimi-go/internal/usecase"
)

var bookLister usecase.BookLister

func main() {
	err := bookLister.ListBooks()
	if err != nil {
		fmt.Println(err)
	}
}
