package natsutils

import (
	"context"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/goutils/cli"
	"testing"
)

func TestNewFromConfig(t *testing.T) {
	assert := assert.New(t)

	config := cli.NewNatsConfig([]string{"nats://localhost:4222"}, "", "")

	assert.Equal(&cli.NatsConfig{
		NatsUrls: cli.NatsUrls{Servers: []string{"nats://localhost:4222"}},
		NatsUser: cli.NatsUser{User: ""},
		NatsPass: cli.NatsPass{Pass: ""},
	}, config)

	natsUtils, err := NewFromConfig(config)
	assert.NoError(err)
	defer natsUtils.Close()

	assert.NotNil(natsUtils.Nats())
	assert.NotNil(natsUtils.Jetstream())

	// use jetstream context to publish
	pub, err := natsUtils.Jetstream().Publish(context.TODO(), "test.test", []byte("test payload"), jetstream.WithExpectStream("TEST"))
	assert.NoError(err)
	assert.NotNil(pub)

	// use natsutils PublishOnStream
	pub, err = natsUtils.PublishOnStream(context.TODO(), "TEST", "test.test", []byte("test payload 2"))
	assert.NoError(err)
	assert.NotNil(pub)

}
