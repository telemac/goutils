package natsservice

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/logger"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/task"
	"time"
)

// NatsServiceRepository provides logging nats cloud events transport for multiple services
type NatsServiceRepository struct {
	task.RunnerRepository
	natsevents.Transporter
	name        string
	natsServers string
	logLevel    string
	logger      *logrus.Entry
	transport   natsevents.Transport // cloudEvent transport, allows to send and receive cloud events over nats
}

// NewNatsServiceRepository creates a nats service repository
func NewNatsServiceRepository(name string, natsServers string, logLevel string) (*NatsServiceRepository, error) {
	nsr := &NatsServiceRepository{name: name, natsServers: natsServers, logLevel: logLevel}

	// create logger
	nsr.logger = logger.New(logLevel, logrus.Fields{
		"service": name,
	})

	// Create nats transport for cloud events
	var err error
	nsr.transport, err = natsevents.NewNatsTransport(natsServers)
	if err != nil {
		return &NatsServiceRepository{}, fmt.Errorf("create nate transport : %w", err)
	}

	return nsr, nil
}

// Close closes flushes events on nats if timeout > 0 and closes nats connection
func (nsr *NatsServiceRepository) Close(timeout time.Duration) error {
	if timeout > 0 {
		_ = nsr.transport.Flush(timeout)
	}
	return nsr.transport.Close()
}

// Start creates a go routine and runs natsSvc Run function
func (nsr *NatsServiceRepository) Start(ctx context.Context, natsSvc NatsServiceIntf, params ...interface{}) *task.Task {
	natsSvc.SetLogger(nsr.logger)
	natsSvc.SetTransport(nsr.transport)
	return nsr.RunnerRepository.Start(ctx, natsSvc, params...)
}

// Logger returns the logger for the service repository
func (nsr *NatsServiceRepository) Logger() *logrus.Entry {
	return nsr.logger.WithField("service", nsr.Name())
}

// Transport returns the transport for the service repository
func (nsr *NatsServiceRepository) Transport() natsevents.Transport {
	return nsr.transport
}

// Name returns the name given in NewNatsTransport
func (nsr *NatsServiceRepository) Name() string {
	return nsr.name
}
