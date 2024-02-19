package natsutils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/megalarm/pkg/cli"
	"testing"
	"time"
)

func TestNatsUtils_CreateConsumer(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	config := cli.NewNatsConfig([]string{"nats://localhost:4222"}, "", "")

	natsUtils, err := NewFromConfig(config)
	assert.NoError(err)
	defer natsUtils.Close()

	assert.NotNil(natsUtils.Nats())
	assert.NotNil(natsUtils.Jetstream())

	// use natsutils PublishOnStream
	var c struct {
		N     int
		Total int
	}
	c.Total = 10
	for i := 0; i < c.Total; i++ {
		c.N = i + 1
		_, err = natsUtils.PublishOnStream(context.TODO(), "TEST", "test.consumer", ToPayload(c))
		assert.NoError(err)
	}

	consumer, err := natsUtils.CreateConsumer(ctx, ConsumerConfig{
		StreamName:     "TEST",
		ConsumerName:   "TEST-DURABLE-CONSUMER",
		Description:    "test durable consumer, can be removed",
		FilterSubjects: []string{"test.consumer"},
		Durable:        true,
	})
	assert.NoError(err)

	//time.Sleep(time.Second)

	moreMessages := true

	for moreMessages {

		startTime := time.Now()

		msgs, err := consumer.Fetch(5, 100*time.Millisecond)
		assert.NoError(err)
		if err != nil {
			moreMessages = false
		}

		for moreMessages {
			select {
			case msg, opened := <-msgs:
				moreMessages = opened
				if msg == nil {
					moreMessages = false
					break
				}
				err = msg.Ack()
				assert.NoError(err)
				fmt.Printf("received %q\n", msg.Subject())
			}
		}
		duration := time.Since(startTime)
		fmt.Printf("duration = %s\n", duration)

		/* OK START*/
		//var receivesMessages int
		//for msg := range msgs {
		//	receivesMessages++
		//	fmt.Printf("received %q\n", msg.Subject())
		//	err = msg.Ack()
		//	assert.NoError(err)
		//	if err != nil {
		//		break
		//	}
		//}
		//if receivesMessages == 0 {
		//	moreMessages = false
		//}
		/* OK END*/

		//for moreMessages {
		//	msg, chanOpened := GetMessage(fetchResult)
		//	_ = chanOpened
		//	if msg != nil {
		//		fmt.Printf("received %q\n", msg.Subject())
		//		err = msg.Ack()
		//		assert.NoError(err)
		//		if err != nil {
		//			moreMessages = false
		//		}
		//	}
		//}

	}

	fmt.Println("the end")

}
