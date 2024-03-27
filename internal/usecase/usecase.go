package usecase

import "github.com/tp86/legimi-go/internal/model"

type BookLister interface {
	ListBooks() ([]model.BookMetadata, error)
}

type BookDownloader interface {
	DownloadBooks([]uint64) error
}
