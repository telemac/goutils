package main

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/events/com.plugis/shell"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("remote-access", "https://nats1.plugis.com", "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("remote-access service starting")

	// start heartbeat service
	servicesRepository.Start(ctx, &heartbeat.HeartbeatService{
		Period:       56,
		RandomPeriod: 3,
	})

	// start shell service
	servicesRepository.Start(ctx, &shell.ShellService{})

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("tempest service ending")
}
