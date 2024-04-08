package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var Options = struct {
	Login             string
	Password          string
	ConfigurationFile string
}{}

func configureFlags() {
	flag.Usage = usage
	flag.StringVar(&Options.ConfigurationFile, "config", "default/location", "config file location")
	flag.StringVar(&Options.Login, "login", "", "Legimi login")
	flag.StringVar(&Options.Password, "password", "", "Legimi password")
}

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
	for _, command := range Commands {
		fmt.Fprintf(&commandsUsage, "  %s %s\n\t%s\n", command.name, command.args, command.description)
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
