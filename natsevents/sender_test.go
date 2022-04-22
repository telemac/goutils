package natsevents

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewCloudEventNats(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const ITERATIONS = 100

	var counter int32
	type Count struct {
		N int
	}

	ConnErrHandler := func(conn *nats.Conn, err error) {
		log.WithError(err).Infoln("nats disconnected")
	}

	ConnHandler := func(conn *nats.Conn) {
		log.Infoln("nats connected")
	}

	opts := cenats.NatsOptions(nats.DisconnectErrHandler(ConnErrHandler), nats.ReconnectHandler(ConnHandler))

	protocol, err := cenats.NewProtocol("https://nats1.plugis.com", "ce.test", "ce.test", opts)
	assert.NoError(err)

	defer protocol.Close(ctx)
	_ = protocol

	c, err := cloudevents.NewClient(protocol)
	assert.NoError(err)

	receiver := func(ctx context.Context, event event.Event) {
		switch event.Type() {
		case "com.drone-box.sample.count":
			var count Count
			err := event.DataAs(&count)
			if err != nil {
				log.WithError(err).Error("event.DataAs")
			}
			atomic.AddInt32(&counter, 1)
			//log.WithField("count", count.N).Info("receiver")
		}
	}

	go func() {
		c.StartReceiver(ctx, receiver)
		log.Info("StartReceiver exited")
	}()

	// wait until connected
	//time.Sleep(time.Second * 2)

	// Create an Event.
	event := cloudevents.NewEvent()
	event.SetSource("example/uri")
	event.SetType("com.drone-box.sample.count")

	var count Count

	// Send that Event.
	for i := 0; i < ITERATIONS; i++ {
		count.N = i
		event.SetData(cloudevents.ApplicationJSON, count)
		event.SetID(uuid.New().String())
		event.SetTime(time.Now())

		result := c.Send(ctx, event)

		if result == nil {
			fmt.Printf("+")
		} else {
			fmt.Printf("!")
		}

		assert.Equal(false, cloudevents.IsUndelivered(result))
		time.Sleep(time.Millisecond * 5)
	}
	time.Sleep(time.Second * 5)
	cancel()
	<-ctx.Done()
	time.Sleep(time.Second * 2)

	assert.Equal(int32(ITERATIONS), atomic.LoadInt32(&counter))
}
