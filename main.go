package main

import (
	"fmt"

	"github.com/tp86/legimi-go/internal/commands"
)

func main() {
	command, err := commands.ParseCommandLine()
	if err != nil {
		fmt.Printf("Error parsing command line: %v\n", err)
		return
	}
	err = command.Run()
	if err != nil {
		fmt.Printf("Error while executing command: %v\n", err)
		return
	}
}
