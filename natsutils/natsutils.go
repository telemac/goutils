package natsutils

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/telemac/goutils/cli"
	"os"
	"path"
)

// NatsUtils provices utilities for nats/jetstream operations
type NatsUtils struct {
	nc         *nats.Conn          // nats connection
	js         jetstream.JetStream // jetstream context
	NatsConfig *cli.NatsConfig
}

// NewFromConnection creates a NatsUtils from an existing nats connection
func NewFromConnection(nc *nats.Conn) (*NatsUtils, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return &NatsUtils{
		nc: nc,
		js: js,
	}, nil
}

// NewFromConfig createx a NatsUtils for a cli.NatsConfig and connects to nats
func NewFromConfig(config *cli.NatsConfig) (*NatsUtils, error) {
	connectStr, err := config.NatsConnectString(0)
	if err != nil {
		return nil, err
	}

	var options []nats.Option
	if true {
		host, _ := os.Hostname()
		exe, _ := os.Executable()
		name := host + "/" + path.Base(exe)
		options = append(options, nats.Name(name))
	}

	nc, err := nats.Connect(connectStr, options...)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return &NatsUtils{
		nc:         nc,
		js:         js,
		NatsConfig: config,
	}, nil

}

// Close closes the nats connection if established
func (natsUtils *NatsUtils) Close() {
	if natsUtils.nc != nil {
		natsUtils.nc.Flush()
		natsUtils.nc.Close()
	}
}

// Nats returns the nats connection
func (natsUtils *NatsUtils) Nats() *nats.Conn {
	return natsUtils.nc
}

// Jetstream returns the jetstream context
func (natsUtils *NatsUtils) Jetstream() jetstream.JetStream {
	return natsUtils.js
}

// PublishOnStream publishes a message on a stream
func (natsUtils *NatsUtils) PublishOnStream(ctx context.Context, stream, subject string, data []byte) (*jetstream.PubAck, error) {
	if natsUtils.js == nil {
		return nil, errors.New("no jetstream context")
	}
	return natsUtils.js.Publish(ctx, subject, data, jetstream.WithExpectStream(stream))
}
