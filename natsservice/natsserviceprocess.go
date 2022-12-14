package natsservice

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/telemac/goutils/task"
	"time"
)

type NatsServiceProcess struct {
	config                NatsServiceProcessConfig
	natsServiceRepository *NatsServiceRepository
	ctx                   context.Context
	cancel                context.CancelFunc
}

type NatsServiceProcessConfig struct {
	ServiceName   string
	NatsServers   string
	LogLevel      string
	CancelTimeout time.Duration
}

func NewServiceProcess(config NatsServiceProcessConfig) (*NatsServiceProcess, error) {
	serviceProcess := &NatsServiceProcess{config: config}

	serviceProcess.ctx, serviceProcess.cancel = task.NewCancellableContext(config.CancelTimeout)

	// create service repository
	var err error
	serviceProcess.natsServiceRepository, err = NewNatsServiceRepository(config.ServiceName, config.NatsServers, config.LogLevel)
	if err != nil {
		log.WithError(err).Error("create nats service repository")
		return nil, err
	}

	serviceProcess.natsServiceRepository.Logger().Info("process started")
	return serviceProcess, nil
}

func (serviceprocess *NatsServiceProcess) Close(timeout time.Duration) error {
	serviceprocess.natsServiceRepository.Logger().Debug("process stoppeing")
	if serviceprocess.cancel != nil {
		serviceprocess.cancel()
	}
	err := serviceprocess.natsServiceRepository.Close(timeout)
	if err != nil {
		serviceprocess.natsServiceRepository.Logger().WithError(err).Error("process stopped")
	} else {
		serviceprocess.natsServiceRepository.Logger().Info("process stopped")
	}
	return err
}

func (serviceprocess *NatsServiceProcess) Start(natsSvc NatsServiceIntf, params ...interface{}) *task.Task {
	return serviceprocess.natsServiceRepository.Start(serviceprocess.ctx, natsSvc, params...)
}

func (serviceprocess *NatsServiceProcess) WaitUntilAllDone() {
	serviceprocess.natsServiceRepository.WaitUntilAllDone()
}

func (serviceprocess *NatsServiceProcess) Logger() *log.Entry {
	return serviceprocess.natsServiceRepository.Logger()
}