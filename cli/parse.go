package cli

import "github.com/alecthomas/kong"

// Parse parses the command line interface and returns the selected command
func Parse(cli any) string {
	ctx := kong.Parse(cli, kong.ConfigureHelp(kong.HelpOptions{Compact: true}))
	command := ctx.Command()
	return command
}
