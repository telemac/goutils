package natsutils

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

type JetstreamConsumer struct {
	config   JetstreamConsumerConfig // config parameters
	nc       *nats.Conn              // nats connection
	js       jetstream.JetStream     // jetstream context
	stream   jetstream.Stream        // jetstream stream
	consumer jetstream.Consumer      // jetstream consumer
}

type JetstreamConsumerConfig struct {
	StreamName       string   // stream from which to consume
	ConsumerName     string   // consumer name
	Description      string   // consumer description
	Durable          bool     // is a durable consumer
	FilterSubjects   []string // subjects received by the consumer
	JsConsumerConfig *jetstream.ConsumerConfig
}

func NewJetstreamConsumer(ctx context.Context, nc *nats.Conn, config JetstreamConsumerConfig) (*JetstreamConsumer, error) {
	var err error
	jetstreamConsumer := &JetstreamConsumer{config: config, nc: nc}

	jetstreamConsumer.js, err = jetstream.New(nc)
	if err != nil {
		return jetstreamConsumer, fmt.Errorf("create jetstream instance : %w", err)
	}

	jetstreamConsumer.stream, err = jetstreamConsumer.js.Stream(ctx, jetstreamConsumer.config.StreamName)
	if err != nil {
		return jetstreamConsumer, fmt.Errorf("get stream : %w", err)
	}

	var consumerConfig jetstream.ConsumerConfig
	if jetstreamConsumer.config.JsConsumerConfig != nil {
		consumerConfig = *jetstreamConsumer.config.JsConsumerConfig
	} else {
		consumerConfig.Name = jetstreamConsumer.config.ConsumerName
		consumerConfig.Description = jetstreamConsumer.config.Description
		consumerConfig.AckPolicy = jetstream.AckExplicitPolicy
		consumerConfig.FilterSubjects = jetstreamConsumer.config.FilterSubjects
		if jetstreamConsumer.config.Durable {
			consumerConfig.Durable = "TEST-DURABLE"
		} else {
			// TODO : adapt parameters for ephemeral consumer
		}
	}

	jetstreamConsumer.consumer, err = jetstreamConsumer.stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		return jetstreamConsumer, fmt.Errorf("create consumer : %w", err)
	}

	return jetstreamConsumer, nil
}

// Fetch returns a channel of messages
func (jetstreamConsumer *JetstreamConsumer) Fetch(batch int, maxWait time.Duration) (jetstream.MessageBatch, error) {
	return jetstreamConsumer.consumer.Fetch(batch, jetstream.FetchMaxWait(maxWait))
}
