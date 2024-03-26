package service

import "github.com/tp86/legimi-go/internal/usecase"

type Session interface {
	GetSession() (string, error)
}

type Account interface {
	GetCredentials() (string, string)
	GetKindleId() uint64
}

type Book interface {
	usecase.BookLister
	usecase.BookDownloader
}
