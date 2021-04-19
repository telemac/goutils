package natsevents

import (
	"context"
	"crypto/tls"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/google/uuid"
	"github.com/jimlawless/whereami"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	net2 "github.com/telemac/goutils/net"
	"net"
	"net/url"
	"reflect"
	"strings"
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

func SNI(serverName string) nats.Option {
	return func(o *nats.Options) error {
		if o.TLSConfig == nil {
			o.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: serverName,
			}
		}
		o.Secure = false
		return nil
	}
}

type SNIDialer struct {
	ServerName string
}

func (d *SNIDialer) Dial(network, address string) (net.Conn, error) {
	conn, err := tls.Dial(network, address, &tls.Config{
		ServerName:         d.ServerName,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	})
	return conn, err
}

// NewNatsTransport connects to nats and creates a NatsTransport
func NewNatsTransport(server string) (*NatsTransport, error) {
	var err error
	transport := &NatsTransport{
		servers:       server,
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

		//nats.Name("cloud1.idronebox.com"),

	}

	// get host name
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	hostname := u.Hostname()

	// TODO : find a solution for SNI with multiple servers, may be it must be implemented in nats
	if strings.Contains(server, "nats://") && strings.Contains(server, ":443") {
		opts = append(opts, nats.SetCustomDialer(&SNIDialer{ServerName: hostname}), SNI(hostname))
	}

	transport.nc, err = nats.Connect(server, opts...)

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
func (t *NatsTransport) Send(ctx context.Context, event *event.Event, topic string) error {
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
func (t *NatsTransport) Request(ctx context.Context, event *event.Event, topic string, timeout time.Duration) (*event.Event, error) {
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
	// receive a valid cloud event, or an enpty payload
	if len(msg.Data) == 0 {
		return nil, nil
	}
	var responseEvent = new(cloudevents.Event)
	err = responseEvent.UnmarshalJSON(msg.Data)
	if err != nil {
		return nil, err
	}

	return responseEvent, nil
}

func (t *NatsTransport) processNatsMessage(cloudEventHandler CloudEventHandler, msg *nats.Msg) {
	var (
		requestEvent  = new(cloudevents.Event)
		responseEvent *cloudevents.Event
		handlerErr    error
	)
	err := requestEvent.UnmarshalJSON(msg.Data)
	if err != nil {
		log.WithError(err).WithField("payload", string(msg.Data)).Warn("decode cloudevent payload")
	}
	if err != nil {
		// call event handler with raw payload
		responseEvent, handlerErr = cloudEventHandler(msg.Subject, nil, msg.Data, err)
	} else {
		// call with decoded event
		responseEvent, handlerErr = cloudEventHandler(msg.Subject, requestEvent, nil, nil)
	}
	if handlerErr != nil {
		// TODO : report this error
		log.WithError(handlerErr).WithField("request", requestEvent).Warn("event handler error")
	}

	if msg.Reply != "" && (responseEvent != nil || handlerErr != nil) {
		// the reply of a cloudevent request must be a valid cloud event.

		if responseEvent == nil {
			// handle response with no event
			responseEvent = t.NewEvent("", requestEvent.Type()+".response", nil)
			EventFillDefaults(responseEvent)
		}

		responseEvent.SetExtension("responsefor", requestEvent.ID())
		if handlerErr != nil {
			responseEvent.SetExtension("error", handlerErr.Error())
		}
		mac, err := net2.GetMACAddress()
		if err != nil {
			log.WithError(err).Error("get mac address")
		}
		responseEvent.SetExtension("mac", mac)

		var payload []byte
		if responseEvent != nil {
			payload, err = responseEvent.MarshalJSON()
			if err != nil {
				// TODO : handle json format errors here
				log.WithError(err).Error("MarshalJSON cloudevent response")

				err = responseEvent.Validate()
				if err != nil {
					log.WithError(err).Error("response event validation")
				}

				responseEvent.SetData(cloudevents.ApplicationJSON, err.Error())
				payload, err = responseEvent.MarshalJSON()
				if err != nil {
					log.WithError(err).Error("MarshalJSON cloudevent response error")
				}
			}
		}
		err = msg.Respond(payload)
		if err != nil {
			// TODO : find a way to report errors in onNatsMessage fn
		}
	}
}

// onNatsMessage is called on each incoming nats message
func (t *NatsTransport) onNatsMessage(msg *nats.Msg) {
	log.WithFields(log.Fields{
		"topic":   msg.Subject,
		"payload": string(msg.Data),
		"reply":   msg.Reply,
	}).Debug("received nats message")

	t.mutex.RLock()
	s, ok := t.subscriptions[msg.Sub.Subject]
	t.mutex.RUnlock()

	if ok {
		go t.processNatsMessage(s.cloudEventHandler, msg)
	} else {
		// message not corresponding a previous subscription, notify ?
		log.WithField("subject", msg.Sub.Subject).Warn("message Subject not subscribed")
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

// NewEvent creates a cloud event given minimal parameters
func (t *NatsTransport) NewEvent(eventPrefix, eventType string, obj interface{}) *event.Event {
	e := event.New(event.CloudEventsVersionV1)
	if eventType != "" {
		e.SetType(eventPrefix + eventType)
	} else {
		objType := reflect.TypeOf(obj).String()
		if strings.HasPrefix(objType, "*") {
			objType = objType[1:]
		}
		e.SetType(eventPrefix + objType)
	}
	e.SetData(event.ApplicationJSON, obj)
	e.SetTime(time.Now())
	e.SetID(uuid.NewString())
	return &e
}

// EventFillDefaults fills the required field with default values if not already set
func EventFillDefaults(e *event.Event) {
	if e.SpecVersion() == "" {
		e.SetSpecVersion(event.CloudEventsVersionV1)
	}
	if e.ID() == "" {
		e.SetID(uuid.NewString())
	}
	if types.IsZero(e.Time()) {
		e.SetTime(time.Now())
	}
	if e.Source() == "" {
		e.SetSource(whereami.WhereAmI(2))
	}
}
