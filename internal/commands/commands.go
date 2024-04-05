package commands

import (
	"fmt"
	"strings"

	"github.com/tp86/legimi-go/internal/options"
	"github.com/tp86/legimi-go/internal/service"
	"github.com/tp86/legimi-go/internal/usecase"
)

type Command func() error

var noneCommand Command = func() error {
	return fmt.Errorf("this should never be called")
}

var commandNames = []string{
	"list",
	"download",
}

func ParseCommandLine() (Command, error) {
	args := options.ParseArgs()
	if len(args) < 1 {
		return noneCommand, fmt.Errorf("expected one of commands: %s\n", strings.Join(commandNames, ", "))
	}
	commandName := args[0]
	switch commandName {
	case "list":
		return listBooks, nil
	case "download":
		return downloadBooks(args[1:]), nil
	default:
		return noneCommand, fmt.Errorf("unsupported command: %s", commandName)
	}
}

var (
	BookLister        usecase.BookLister
	BookListPresenter service.BookListPresenter
	BookDownloader    usecase.BookDownloader
)

func listBooks() error {
	bookList, err := BookLister.ListBooks()
	if err != nil {
		return err
	}
	BookListPresenter.Present(bookList)
	return nil
}

func downloadBooks(ids []string) Command {
	return func() error {
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
		return BookDownloader.DownloadBooks(bookIds)
	}
}
