package cli

// LogLevel is the log level
type LogLevel struct {
	Log string `help:"log level (trace|debug|info|error)" short:"l" default:"info"`
}
