package config

import "flag"

type EventSenderConfig struct {
	Server    string
	LogLevel  string
	EventType string
	EventData string
	Topic     string
	Request   bool
	Timeout   int
}

func (c *EventSenderConfig) Parse() {
	flag.StringVar(&c.Server, "server", "https://nats1.plugis.com", "nats server")
	flag.StringVar(&c.LogLevel, "log", "info", "log level (trace|debug|info|warn|error|fatal)")
	flag.StringVar(&c.EventType, "type", "com.plugis.browser.open", "cloud event type")
	flag.StringVar(&c.EventData, "data", `{"url":"https://google.fr"}`, "cloud event data in json format")
	flag.StringVar(&c.Topic, "topic", `com.plugis.browser`, "topic to send the event")
	flag.BoolVar(&c.Request, "request", false, "send request")
	flag.IntVar(&c.Timeout, "timeout", 60, "timeout")
	flag.Parse()
}
