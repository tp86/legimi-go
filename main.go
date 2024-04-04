package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tp86/legimi-go/internal/service"
	"github.com/tp86/legimi-go/internal/usecase"
)

var (
	bookLister        usecase.BookLister
	bookListPresenter service.BookListPresenter
	bookDownloader    usecase.BookDownloader
)

type command struct {
	Exec func() error
}

var commandNames = []string{
	"list",
	"download",
}

func main() {
	command, err := parseCommand()
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}
	err = command.Exec()
	if err != nil {
		fmt.Printf("Error while executing command: %v\n", err)
		return
	}
}

func parseCommand() (command, error) {
	// configFile := flag.String("config", "default/location", "config file location")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return command{}, fmt.Errorf("expected one of commands: %s\n", strings.Join(commandNames, ", "))
	}
	commandName := args[0]
	switch commandName {
	case "list":
		return command{Exec: listBooks}, nil
	case "download":
		return command{Exec: func() error {
			return downloadBooks(args[1:])
		}}, nil
	default:
		return command{}, fmt.Errorf("unsupported command: %s", commandName)
	}
}

func listBooks() error {
	bookList, err := bookLister.ListBooks()
	if err != nil {
		return err
	}
	bookListPresenter.Present(bookList)
	return nil
}

func downloadBooks(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("no book id provided")
	}
	bookIds := make([]uint64, len(ids))
	for i, id := range ids {
		_, err := fmt.Sscanf(id, "%d", &bookIds[i])
		if err != nil {
			return err
		}
	}
	return bookDownloader.DownloadBooks(bookIds)
}
