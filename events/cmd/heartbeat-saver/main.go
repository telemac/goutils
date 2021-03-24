package main

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	config, err := natsservice.LoadConfig("./heartbeat-saver.yml")
	if err != nil {
		logrus.WithError(err).Fatal("open config file")
	}

	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("heartbeat-saver", "https://nats1.plugis.com", "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("heartbeat-saver service starting")

	// start heartbeat saver
	servicesRepository.Start(ctx, heartbeat.NewHeartbeatSaver(config.Mysql))

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("heartbeat-saver service ending")
}
