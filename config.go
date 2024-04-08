package main

import (
	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/commands"
	ar "github.com/tp86/legimi-go/internal/repository/account"
	as "github.com/tp86/legimi-go/internal/service/account"
	"github.com/tp86/legimi-go/internal/service/book"
	"github.com/tp86/legimi-go/internal/service/presenter"
	"github.com/tp86/legimi-go/internal/service/session"
)

func configure() {
	accountRepository := ar.GetMemoryRepository()
	accountService := as.DefaultService(accountRepository, commands.Options.Login, commands.Options.Password)
	apiClient := api.GetClient()
	sessionService := session.DefaultService(accountService, apiClient)
	bookDownloadPresenter := presenter.DefaultBookDownloadPresenter()
	bookService := book.DefaultService(sessionService, apiClient, bookDownloadPresenter)
	commands.BookLister = bookService
	commands.BookDownloader = bookService
	commands.BookListPresenter = presenter.DefaultBookListPresenter()
}
