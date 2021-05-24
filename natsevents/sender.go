package natsevents

import (
	"context"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
)

// CloudEventSender Sends cloud events
type CloudEventSender interface {
	Flush(timeout time.Duration) error
	Send(ctx context.Context, event *event.Event, topic string) error
	Request(ctx context.Context, event *event.Event, topic string, timeout time.Duration) (*event.Event, error)
}
