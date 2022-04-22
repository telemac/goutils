package commandline

import (
	"flag"
)

// CommandLine holds the minimal common command line parameters
type CommandLine struct {
	Log          string // log level (trace|debug|info|warn|error|fatal)
	ConfigFolder string // folder containing config files, default ./
}

// Parse parses the command line parameters and returns a CommandLine instance
func Parse() (commandLine CommandLine) {
	flag.StringVar(&commandLine.Log, "log", "info", "log level (trace|debug|info|warn|error|fatal)")
	flag.StringVar(&commandLine.ConfigFolder, "config", "./", "config file directory path")
	flag.Parse()
	return
}
