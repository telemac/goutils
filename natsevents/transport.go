package natsevents

import (
	"io"
)

type Transport interface {
	CloudEventSender
	CloudEventReceiver
	io.Closer
}
