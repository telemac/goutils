package natsutils

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewJetstreamConsumer(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	natsConnectString := "nats://localhost:4222"

	nc, err := nats.Connect(natsConnectString)
	assert.NoError(err)
	defer nc.Close()

	jsConsumer, err := NewJetstreamConsumer(ctx, nc, JetstreamConsumerConfig{
		StreamName:   "WATCHCOM",
		ConsumerName: "TEST-DURABLE",
		Description:  "durable test consumer test",
		Durable:      true,
		FilterSubjects: []string{
			"watchcom.watchdog",
			"watchcom.alarm",
		},
	})
	assert.NoError(err)

	for {
		fetchResult, err := jsConsumer.Fetch(1, 10*time.Second)
		fetchResult.Error()
		assert.NoError(err)
		if err != nil {
			return
		}
		for msg := range fetchResult.Messages() {
			fmt.Printf("received %q\n", msg.Subject())

			var ce event.Event // cloud event
			err = ce.UnmarshalJSON(msg.Data())
			assert.NoError(err)

			eventType := ce.Type()

			fmt.Printf("ce type = %s\n", eventType)

			switch eventType {
			case "com.megalarm.watchcom.watchdog":
				//var wf model.WatchcomFrame
				//err = ce.DataAs(&wf)
				//assert.NoError(err)
			}

			err = msg.Ack()
			assert.NoError(err)
			if err != nil {
				return
			}
		}
	}

}
