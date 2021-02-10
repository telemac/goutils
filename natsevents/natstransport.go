package natsevents

import (
	"context"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/jimlawless/whereami"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// NatsTransport allows to send and receive CloudEvents on nats
type NatsTransport struct {
	CloudEventSender
	CloudEventReceiver
	servers       string
	nc            *nats.Conn
	mutex         sync.RWMutex
	subscriptions map[string]*natsSubscription // subscriptions done by RegisterHandler, indexed by subscription.Subject
}

// natsSubscription holds the nats.Subscription and callback
type natsSubscription struct {
	subscription      *nats.Subscription
	cloudEventHandler CloudEventHandler
}

// Close closes the underlying nats connection
func (t *NatsTransport) Close() error {
	if t.nc == nil {
		return errors.New("nc not initialized by nats.Connect")
	}

	// TODO : may be t.Flush here
	// t.Flush(time.Second * 5)

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, natsSubscription := range t.subscriptions {
		natsSubscription.subscription.Drain()
		// wait until all pending messages received
		for {
			msgs, _, err := natsSubscription.subscription.Pending()
			if err == nil && msgs > 0 {
				time.Sleep(time.Millisecond * 10)
			} else {
				break
			}
		}
	}

	// TODO : Drain before drain subscriptions ?
	t.nc.Drain()

	t.nc.Close()
	t.nc = nil
	return nil
}

// Connected returns true if nats is connected
func (t *NatsTransport) Connected() bool {
	if t.nc == nil {
		return false
	}
	return t.nc.IsConnected()
}

// NewNatsTransport connects to nats and creates a NatsTransport
func NewNatsTransport(servers string) (*NatsTransport, error) {
	var err error
	transport := &NatsTransport{
		servers:       servers,
		subscriptions: make(map[string]*natsSubscription),
	}

	opts := []nats.Option{
		nats.ReconnectHandler(transport.reconnectedCB),
		nats.DisconnectErrHandler(transport.disconnectedErrCB),
		nats.ClosedHandler(transport.closedHandler),
		nats.DiscoveredServersHandler(transport.discoveredServersCB),
		nats.ErrorHandler(transport.errorHandler),
		nats.MaxReconnects(-1),
		nats.DrainTimeout(30 * time.Second),
	}

	transport.nc, err = nats.Connect(servers, opts...)

	return transport, err
}

func (t *NatsTransport) reconnectedCB(conn *nats.Conn) {
	log.Infoln("nats reconnectedCB")
}

func (t *NatsTransport) disconnectedErrCB(conn *nats.Conn, err error) {
	if err != nil {
		log.WithError(err).Warnln("nats ConnErrHandler")
	} else {
		log.Infoln("nats ConnErrHandler")
	}
}

func (t *NatsTransport) closedHandler(conn *nats.Conn) {
	log.Infoln("nats closedHandler")
}

func (t *NatsTransport) discoveredServersCB(conn *nats.Conn) {
	log.Infoln("nats discoveredServersCB")
}

func (t *NatsTransport) errorHandler(conn *nats.Conn, subscription *nats.Subscription, err error) {
	log.WithError(err).WithField("subscription", subscription).Error("nats errorHandler")
}

// Flush flushes nats pending writes
func (t *NatsTransport) Flush(timeout time.Duration) error {
	return t.nc.FlushTimeout(timeout)
}

// Send sends the json representation of the event on nats topic
func (t *NatsTransport) Send(ctx context.Context, event event.Event, topic string) error {
	if event.Source() == "" {
		event.SetSource(whereami.WhereAmI(2))
	}
	payload, err := event.MarshalJSON()
	if err != nil {
		return err
	}
	err = t.nc.Publish(topic, payload)
	if errors.Is(err, nats.ErrReconnectBufExceeded) {
		t.Flush(time.Second * 5)
		err = t.nc.Publish(topic, payload)
	}
	return err
}

// Request sends a request and waits the response for timeout
func (t *NatsTransport) Request(ctx context.Context, event event.Event, topic string, timeout time.Duration) (*event.Event, error) {
	if event.Source() == "" {
		event.SetSource(whereami.WhereAmI(2))
	}

	payload, err := event.MarshalJSON()
	if err != nil {
		return nil, err
	}
	t.Flush(time.Millisecond * 10)
	msg, err := t.nc.Request(topic, payload, timeout)
	if err != nil {
		return nil, err
	}
	// TODO : check enpty payload response
	var responseEvent = new(cloudevents.Event)
	err = responseEvent.UnmarshalJSON(msg.Data)
	if err != nil {
		return nil, err
	}

	return responseEvent, nil
}

// onNatsMessage is called on each incoming nats message
func (t *NatsTransport) onNatsMessage(msg *nats.Msg) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	s, ok := t.subscriptions[msg.Sub.Subject]
	if ok {
		var event = new(cloudevents.Event)
		err := event.UnmarshalJSON(msg.Data)
		if err != nil {
			// call event handler with raw payload
			event, err = s.cloudEventHandler(msg.Subject, nil, msg.Data)
		} else {
			// call with decoded event
			event, err = s.cloudEventHandler(msg.Subject, event, nil)
		}
		// TODO : check if it is a request response
		if msg.Reply != "" {
			payload, err := event.MarshalJSON()
			if err != nil {
				// TODO : handle json format errors here
			}
			msg.Respond(payload)
		}
	} else {
		// TODO : message not corresponding a previous subscription, notify ?
	}
}

func (t *NatsTransport) RegisterHandler(eventHandler CloudEventHandler, topic string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// check if subscription already done
	_, ok := t.subscriptions[topic]
	if ok {
		return nil
	}

	subscription, err := t.nc.Subscribe(topic, t.onNatsMessage)
	if err != nil {
		return err
	}
	natsSub := &natsSubscription{
		subscription:      subscription,
		cloudEventHandler: eventHandler,
	}
	t.subscriptions[subscription.Subject] = natsSub

	return nil
}

var contextKey struct{}

// WithTransport adds the transport as context value
func WithTransport(ctx context.Context, transport *NatsTransport) context.Context {
	return context.WithValue(ctx, contextKey, transport)
}

// FromContext returns the transport set with WithTransport or nil
func FromContext(ctx context.Context) *NatsTransport {
	transport, ok := ctx.Value(contextKey).(*NatsTransport)
	if !ok {
		return nil
	}
	return transport
}
