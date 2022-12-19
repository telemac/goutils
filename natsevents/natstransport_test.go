package natsevents

import (
	"context"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type Count struct {
	N int
}

func TestNewNatsTransport(t *testing.T) {
	var lastCount = 0
	const ITERATIONS = 2000

	logrus.SetLevel(logrus.TraceLevel)
	assert := assert.New(t)
	//transport, err := NewNatsTransport("localhost")
	var transport *NatsTransport
	var err error
	transport, err = NewNatsTransport("nats://cloud1.idronebox.com:443")
	assert.NoError(err)
	assert.True(transport.Connected())
	//transport, err = NewNatsTransport("nats://server1.plugis.com:443")
	//assert.NoError(err)
	//assert.True(transport.Connected())

	eventHandler := func(topic string, event *event.Event, payload []byte, err error) (*event.Event, error) {
		//logrus.Printf("eventHandler callback topic = %s event = %s", topic, event)
		if event == nil {
			logrus.Fatal("event is nil in event handler")
			return nil, errors.New("event not defined")
		}
		switch topic {
		case "events.test":
			t := event.Type()
			switch t {
			case "com.plugis.sample.event.natsevents.Count":
				var count Count
				err = event.DataAs(&count)
				assert.NoError(err)
				if err != nil {
					logrus.WithError(err).Fatal("decode event data")
				}
				if lastCount != count.N-1 {
					assert.Equal(count.N-1, lastCount)
					//logrus.Fatal("lastCount")
				}
				lastCount = count.N
				//logrus.Printf("count = %d", count.N)
			default:
				logrus.WithFields(logrus.Fields{
					"type":  t,
					"topic": topic,
				}).Fatal("unknown event type")
			}
		}
		return nil, nil
	}

	err = transport.RegisterHandler(eventHandler, "events.>")
	assert.NoError(err)

	time.Sleep(time.Second * 1)

	// Create an Event.
	event := cloudevents.NewEvent()
	//event.SetSource("pkg/natsevents/natstransport_test.go")
	ctx := context.TODO()

	// send an event
	var count Count
	for i := 1; i <= ITERATIONS; i++ {
		count.N = i
		event.SetType("com.plugis.sample.event." + reflect.TypeOf(count).String())
		event.SetID(uuid.New().String())
		event.SetTime(time.Now())
		event.SetData(cloudevents.ApplicationJSON, count)
		err = transport.Send(ctx, &event, "events.test")
		//assert.NoError(err, "transport.Send")
		if err != nil {
			logrus.WithError(err).Error("transport.Send")
		}
	}

	err = transport.Flush(time.Second * 30)
	assert.NoError(err)

	err = transport.Close()
	assert.NoError(err)
	assert.False(transport.Connected())

	assert.Equal(ITERATIONS, lastCount)

}

func TestNatsTransport_Request(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	assert := assert.New(t)
	//transport, err := NewNatsTransport("localhost")
	transport, err := NewNatsTransport("https://nats1.plugis.com")
	//transport, err := NewNatsTransport("https://demo.nats.io")
	assert.NoError(err)
	assert.True(transport.Connected())

	ctx := context.Background()

	// Create an Event.
	type SumRequest struct {
		A, B int
		Sum  int
	}

	req := SumRequest{A: 10, B: 20}
	ce := cloudevents.NewEvent()

	ce.SetType("com.plugis.sample.request." + reflect.TypeOf(req).String())
	ce.SetID(uuid.New().String())
	ce.SetTime(time.Now())
	ce.SetData(cloudevents.ApplicationJSON, req)
	responseCloudEvent, err := transport.Request(ctx, &ce, "events.request", time.Second)
	assert.True(errors.Is(err, nats.ErrTimeout))

	// add handler
	requestHandler := func(topic string, event *event.Event, payload []byte, err error) (*event.Event, error) {
		//logrus.Printf("eventHandler callback topic = %s ce = %s", topic, ce)
		if event == nil {
			logrus.Fatal("ce is nil in ce handler")
			return nil, errors.New("ce not defined")
		}
		switch topic {
		case "events.request":
			t := event.Type()
			switch t {
			case "com.plugis.sample.request.natsevents.SumRequest":
				var sumReq SumRequest
				err = event.DataAs(&sumReq)
				if err != nil {
					logrus.WithError(err).Fatal("decode ce data")
					return nil, err
				}
				sumReq.Sum = sumReq.A + sumReq.B
				ce.SetData(cloudevents.ApplicationJSON, sumReq)
				return &ce, nil
				//logrus.Printf("count = %d", count.N)
			default:
				logrus.WithFields(logrus.Fields{
					"type":  t,
					"topic": topic,
				}).Fatal("unknown ce type")
			}
		}
		return nil, nil
	}

	err = transport.RegisterHandler(requestHandler, "events.request")
	assert.NoError(err)

	ce.SetSource("")
	responseCloudEvent, err = transport.Request(ctx, &ce, "events.request", time.Second)
	assert.NoError(err)

	assert.NotNil(responseCloudEvent)

	// get data from cloudEvent
	var sumResp SumRequest
	err = responseCloudEvent.DataAs(&sumResp)
	assert.NoError(err)
	assert.Equal(SumRequest{A: 10, B: 20, Sum: 30}, sumResp)

	transport.Close()
}
