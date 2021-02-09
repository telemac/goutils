package natsevents

import "github.com/cloudevents/sdk-go/v2/protocol"

type Transport interface {
	CloudEventSender
	CloudEventReceiver
	protocol.Closer
}
