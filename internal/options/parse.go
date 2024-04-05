package options

import (
	"flag"
)

func ParseArgs() []string {
	flag.StringVar(ConfigurationFile, "config", "default/location", "config file location")
	flag.Parse()
	return flag.Args()
}
