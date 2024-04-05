package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/commands"
	"github.com/tp86/legimi-go/internal/options"
	ar "github.com/tp86/legimi-go/internal/repository/account"
	as "github.com/tp86/legimi-go/internal/service/account"
	"github.com/tp86/legimi-go/internal/service/book"
	"github.com/tp86/legimi-go/internal/service/presenter"
	"github.com/tp86/legimi-go/internal/service/session"
)

func flagPrint(what string) {
	fmt.Fprint(flag.CommandLine.Output(), what)
}

func usage() {
	usageHeader := fmt.Sprintf(
		"Alternative downloader of Legimi ebooks for Kindle.\n\n"+
			"Usage:\n  %s [options] <command> [arguments]\n",
		os.Args[0])
	commandsHeader := "Commands are:\n"
	var commandsUsage strings.Builder
	for _, command := range commands.Commands {
		fmt.Fprintf(&commandsUsage, "  %s %s\n\t%s\n", command.Name, command.Args, command.Description)
	}
	optionsHeader := "General options are:\n"
	flagPrint(usageHeader)
	flagPrint("\n")
	flagPrint(commandsHeader)
	flagPrint(commandsUsage.String())
	flagPrint("\n")
	flagPrint(optionsHeader)
	flag.PrintDefaults()
	flagPrint("\n")
}

func init() {
	flag.Usage = usage
	accountRepository := ar.GetMemoryRepository()
	accountService := as.DefaultService(accountRepository, options.Login, options.Password)
	apiClient := api.GetClient()
	sessionService := session.DefaultService(accountService, apiClient)
	bookDownloadPresenter := presenter.DefaultBookDownloadPresenter()
	bookService := book.DefaultService(sessionService, apiClient, bookDownloadPresenter)
	commands.BookLister = bookService
	commands.BookDownloader = bookService
	commands.BookListPresenter = presenter.DefaultBookListPresenter()
}
