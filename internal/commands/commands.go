package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/tp86/legimi-go/internal/options"
	"github.com/tp86/legimi-go/internal/service"
	"github.com/tp86/legimi-go/internal/usecase"
)

type Command struct {
	Name        string
	Args        string
	Description string
	Run         func() error
}

var noneCommand = Command{
	Run: func() error {
		return fmt.Errorf("this should never be called")
	},
}

var (
	Commands = []Command{
		{Name: "list", Run: listBooks, Description: "list books on shelf"},
		{Name: "download", Args: "id ...", Run: downloadBooks, Description: "download book(s) with given id(s)"},
	}
)

func ParseCommandLine() (Command, error) {
	options.Configure()
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		names := make([]string, len(Commands))
		for i, command := range Commands {
			names[i] = command.Name
		}
		return noneCommand, fmt.Errorf("expected one of commands: %s\n", strings.Join(names, ", "))
	}
	commandName := args[0]
	if command, ok := findCommand(commandName); ok {
		return command, nil
	}
	return noneCommand, fmt.Errorf("unsupported command: %s", commandName)
}

func findCommand(name string) (Command, bool) {
	for _, command := range Commands {
		if command.Name == name {
			return command, true
		}
	}
	return noneCommand, false
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

func downloadBooks() error {
	ids := flag.Args()[1:]
	if len(ids) == 0 {
		return fmt.Errorf("no book id provided")
	}
	bookIds := make([]uint64, len(ids))
	for i, id := range ids {
		v, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			err := err.(*strconv.NumError)
			return fmt.Errorf("couldn't parse id '%s': %s", id, err.Err)
		}
		bookIds[i] = v
	}
	return BookDownloader.DownloadBooks(bookIds)
}
