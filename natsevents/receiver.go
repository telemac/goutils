package natsevents

import (
	"github.com/cloudevents/sdk-go/v2/event"
)

// CloudEventHandler : callback to handle CloudEvents
// must return as fast as possible.
// event can be nil, in this case the payload will be set
// if event is not nil (properly decoded) the payload is nil
type CloudEventHandler func(topic string, event *event.Event, payload []byte, err error) (*event.Event, error)

// CloudEventReceiver allows to receive cloud events
type CloudEventReceiver interface {
	RegisterHandler(eventHandler CloudEventHandler, topic string) error
	Respond(*event.Event, error)
}
