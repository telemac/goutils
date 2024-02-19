package cli

// LogLevel is the log level
type TCPPort struct {
	Port int `help:"port" short:"P" min:"0" max:"65535"`
}
