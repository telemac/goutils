package udp

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*100))
	defer cancel()

	err := Broadcast("255.255.255.255:50222", []byte("hello"))
	assert.NoError(err)

	<-ctx.Done()
}

func TestNewUDPBroadcaster(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*1))
	defer cancel()

	udpBroadcaster, err := NewUDPBroadcaster(50222)
	assert.NoError(err)

	payload := []byte("hello from udpBroadcaster")

	err = udpBroadcaster.Broadcast(payload)
	assert.NoError(err)

	datagram, err := udpBroadcaster.Read(time.Millisecond * 100)
	assert.NoError(err)
	assert.Equal(payload, datagram.Payload)
	assert.Contains(datagram.Addr.String(), ":50222")

	datagram, err = udpBroadcaster.Read(time.Millisecond * 100)
	assert.True(errors.Is(err, os.ErrDeadlineExceeded))
	<-ctx.Done()

	udpBroadcaster.Close()
}
