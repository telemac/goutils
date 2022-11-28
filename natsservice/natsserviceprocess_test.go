package natsservice

import (
	"github.com/stretchr/testify/assert"
	"github.com/telemac/goutils/natsevents"
	"golang.org/x/net/context"
	"testing"
	"time"
)

type Svc struct {
	NatsService
}

func (s *Svc) Run(ctx context.Context, params ...interface{}) error {
	ce := natsevents.NewEvent("", "com.plugis.test", "this is a test")
	err := s.transport.Send(ctx, ce, "test")
	if err != nil {
		s.logger.WithError(err).Error("send cloudevent")
		return err
	}
	return s.transport.Flush(time.Second * 3)
}

func TestNewServiceProcess(t *testing.T) {
	assert := assert.New(t)
	serviceProcess, err := NewServiceProcess(NatsServiceProcessConfig{
		ServiceName:   "sms-receiver",
		NatsServers:   "nats://demo.nats.io:4222",
		LogLevel:      "trace",
		CancelTimeout: time.Second * 15,
	})
	assert.NoError(err)
	assert.NotNil(serviceProcess)

	defer func() {
		err := serviceProcess.Close(time.Second * 10)
		assert.NoError(err)
	}()

	svc := new(Svc)
	serviceProcess.Start(svc)

	serviceProcess.WaitUntilAllDone()

}
