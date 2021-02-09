package natsservice

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/goutils/logger"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/task"
	"testing"
	"time"
)

// create a sample service
type SampleNatsService struct {
	*NatsService
}

func NewSampleNatsService() (*SampleNatsService, error) {
	log := logger.New("trace", logrus.Fields{
		"service": "sample nats service",
	})

	// Create nats transport for cloud events
	transport, err := natsevents.NewNatsTransport("https://nats1.plugis.com")
	if err != nil {
		return &SampleNatsService{}, fmt.Errorf("create nate transport : %w", err)
	}

	return &SampleNatsService{NatsService: NewNatsService(log, transport)}, nil
}

func (svc *SampleNatsService) Run(ctx context.Context, params ...interface{}) error {
	svc.logger.Info("service started")
	defer svc.logger.Info("service stopped")
	for !task.IsCancelled(ctx) {
		svc.logger.Info("service tick")
		task.Sleep(ctx, time.Second*1)
	}
	return nil
}

func TestNewNatsService(t *testing.T) {
	ctx,cancel := task.NewCancellableContext(time.Second*10)
	ctx,cancel = context.WithTimeout(ctx,time.Second*5)
	defer cancel()
	assert := assert.New(t)
	sampleService,err := NewSampleNatsService()
	assert.NoError(err)
	err = sampleService.Run(ctx)
	sampleService.transport.Close()
}
