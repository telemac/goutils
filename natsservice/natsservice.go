package natsservice

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/logger"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/task"
)

type NatsServiceIntf interface {
	task.Runnable
	logger.Logger
	natsevents.Transporter
}

type NatsService struct {
	logger    *logrus.Entry
	transport natsevents.Transport // cloudEvent transport, allows to send and receive cloud events over nats
	NatsServiceIntf
}

// implement Logger
func (ns *NatsService) Logger() *logrus.Entry {
	if ns.logger == nil {
		panic("logger not set in NatsService")
	}
	return ns.logger
}

func (ns *NatsService) SetLogger(logger *logrus.Entry) {
	ns.logger = logger
}

// implement Transporter
func (ns *NatsService) Transport() natsevents.Transport {
	if ns.transport == nil {
		if ns.Logger() != nil {
			ns.Logger().Error("transport not set in NatsService")
		}
		panic("transport not set in NatsService")
	}
	return ns.transport
}

func (ns *NatsService) SetTransport(transport natsevents.Transport) {
	ns.transport = transport
}

// NewNatsService creates a nats service
func NewNatsService(logger *logrus.Entry, transport natsevents.Transport) *NatsService {
	return &NatsService{logger: logger, transport: transport}
}
