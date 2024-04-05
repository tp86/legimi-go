package options

import (
	"flag"
)

var (
	configurationFile string
	ConfigurationFile = &configurationFile
)

func Configure() {
	flag.StringVar(ConfigurationFile, "config", "default/location", "config file location")
}
