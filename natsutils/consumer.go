package natsutils

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

type Consumer struct {
	natsUtils  *NatsUtils
	config     ConsumerConfig
	stream     jetstream.Stream
	jsConsumer jetstream.Consumer
}

type ConsumerConfig struct {
	StreamName        string        // stream from which to consume
	ConsumerName      string        // consumer name
	Description       string        // consumer description
	FilterSubjects    []string      // subjects received by the consumer
	Durable           bool          // is a durable consumer
	InactiveThreshold time.Duration // used if durable is true only
}

func (natsUtils *NatsUtils) CreateConsumer(ctx context.Context, config ConsumerConfig) (*Consumer, error) {
	var err error

	consumer := &Consumer{natsUtils: natsUtils, config: config}

	consumer.stream, err = natsUtils.js.Stream(ctx, consumer.config.StreamName)
	if err != nil {
		return consumer, fmt.Errorf("get stream : %w", err)
	}

	jsConsumerConfig := jetstream.ConsumerConfig{
		Name:           consumer.config.ConsumerName,
		Description:    consumer.config.Description,
		AckPolicy:      jetstream.AckExplicitPolicy,
		FilterSubjects: consumer.config.FilterSubjects,
	}
	if consumer.config.Durable {
		jsConsumerConfig.Durable = consumer.config.ConsumerName
	} else {
		// adapt parameters for ephemeral consumer
		if consumer.config.InactiveThreshold.Nanoseconds() == 0 {
			jsConsumerConfig.InactiveThreshold = time.Second * 60 // default 1 mn
		} else {
			jsConsumerConfig.InactiveThreshold = consumer.config.InactiveThreshold
		}
	}

	consumer.jsConsumer, err = consumer.stream.CreateOrUpdateConsumer(ctx, jsConsumerConfig)
	if err != nil {
		return consumer, fmt.Errorf("create consumer : %w", err)
	}

	return consumer, nil
}

// Fetch returns a channel of messages immediately, the channel will be closed in maxWait or batch messages received
func (consumer *Consumer) Fetch(batch int, maxWait time.Duration) (<-chan jetstream.Msg, error) {
	msgs, err := consumer.jsConsumer.Fetch(batch, jetstream.FetchMaxWait(maxWait))
	if err != nil {
		return nil, err
	}
	return msgs.Messages(), nil
}
