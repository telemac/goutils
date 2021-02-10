package natsservice

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/logger"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/task"
)

// NatsServiceRepository provides logging nats cloud events transport for multiple services
type NatsServiceRepository struct {
	task.RunnerRepository
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
	nsr.transport, err = natsevents.NewNatsTransport("https://nats1.plugis.com")
	if err != nil {
		return &NatsServiceRepository{}, fmt.Errorf("create nate transport : %w", err)
	}

	return nsr, nil
}

func (nsr *NatsServiceRepository) Close() error {
	return nsr.transport.Close()
}

func (nsr *NatsServiceRepository) Start(ctx context.Context, natsSvc NatsServiceIntf, params ...interface{}) *task.Task {
	natsSvc.SetLogger(nsr.logger)
	natsSvc.SetTransport(nsr.transport)
	return nsr.RunnerRepository.Start(ctx, natsSvc, params...)
}
