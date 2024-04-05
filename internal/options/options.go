package options

import (
	"flag"
)

var (
	configurationFile string
	ConfigurationFile = &configurationFile

	login string
	Login = &login

	password string
	Password = &password
)

func Configure() {
	flag.StringVar(ConfigurationFile, "config", "default/location", "config file location")
	flag.StringVar(Login, "login", "", "Legimi login")
	flag.StringVar(Password, "password", "", "Legimi password")
}
