package natsevents

import (
	"io"
)

type Transport interface {
	CloudEventSender
	CloudEventReceiver
	io.Closer
}

type Transporter interface {
	Transport() Transport
	SetTransport(transport Transport)
}
