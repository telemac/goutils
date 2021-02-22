package main

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("heartbeat-saver", "https://nats1.plugis.com", "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("heartbeat-saver service starting")

	// start heartbeat saver
	servicesRepository.Start(ctx, &heartbeat.HeartbeatSaver{})

	// start heartbeat sender
	//servicesRepository.Start(ctx, &heartbeat.HeartbeatSender{
	//	Period:       55,
	//	RandomPeriod: 4,
	//})

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("tempest service ending")
}
