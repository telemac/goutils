package natsservice

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/task"
)

type NatsService struct {
	logger *logrus.Entry
	transport natsevents.Transport // cloudEvent transport, allows to send and receive cloud events over nats
	task.Runnable
}

// NewNatsService creates a nats service
func NewNatsService(logger *logrus.Entry, transport natsevents.Transport) *NatsService {
	return &NatsService{logger: logger, transport: transport}
}
