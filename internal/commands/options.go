package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func defaultConfigurationFile() string {
	return os.Getenv("HOME") + "/.config/legimi-go/config.ini"
}

type options struct {
	login             string
	password          string
	configurationFile string
	debugging         bool
}

var Options options

func (o options) GetLogin() string {
	return o.login
}

func (o options) GetPassword() string {
	return o.password
}

func (o options) GetFile() string {
	return o.configurationFile
}

func configureFlags() {
	flag.Usage = usage
	flag.StringVar(&Options.configurationFile, "config", defaultConfigurationFile(), "path to configuration file")
	flag.StringVar(&Options.login, "login", "", "Legimi login")
	flag.StringVar(&Options.password, "password", "", "Legimi password")
	flag.BoolVar(&Options.debugging, "debug", false, "print debugging information")
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
