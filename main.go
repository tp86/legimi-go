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

var commandMap = map[string]func([]string) error{
	"list":     listBooks,
	"download": downloadBooks,
}

func main() {
	configFile := flag.String("config", "default/location", "config file location")
	flag.Parse()
	fmt.Println(*configFile)
	args := flag.Args()
	commands := make([]string, 0, len(commandMap))
	for name := range commandMap {
		commands = append(commands, name)
	}
	if len(args) < 1 {
		fmt.Printf("expected one of commands: %s\n", strings.Join(commands, ", "))
		return
	}
	commandName := args[0]
	command, ok := commandMap[commandName]
	if !ok {
		fmt.Printf("unknown command: %s\n", commandName)
		return
	}
	err := command(args[1:])
	if err != nil {
		fmt.Println(err)
	}
}

func listBooks([]string) error {
	bookList, err := bookLister.ListBooks()
	if err != nil {
		return err
	}
	bookListPresenter.Present(bookList)
	return nil
}

func downloadBooks(ids []string) error {
	if len(ids) < 1 {
		return fmt.Errorf("no book id provided")
	}
	// only first book id supported for now
	var bookId uint64
	n, err := fmt.Sscanf(ids[0], "%d", &bookId)
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("couldn't parse book id")
	}
	return bookDownloader.DownloadBooks([]uint64{uint64(bookId)})
}
