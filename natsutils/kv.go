package natsutils

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go/jetstream"
)

// CreateKeyValue creates a key value store
func (natsUtils *NatsUtils) CreateKeyValue(ctx context.Context, config jetstream.KeyValueConfig) (jetstream.KeyValue, error) {
	if natsUtils.js == nil {
		return nil, errors.New("no jetstream context")
	}
	return natsUtils.js.CreateKeyValue(ctx, config)
}
