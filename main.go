package main

import (
	"fmt"
	"os"

	"github.com/tp86/legimi-go/internal/commands"
)

func main() {
	command, err := commands.ParseCommandLine()
	if err != nil {
		fmt.Printf("Error parsing command line: %v\n", err)
		os.Exit(2)
	}
	configure()
	err = command.Run()
	if err != nil {
		fmt.Printf("Error while executing command: %v\n", err)
		os.Exit(1)
	}
}
