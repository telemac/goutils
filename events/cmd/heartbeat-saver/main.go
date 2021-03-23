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

	dbConfig := heartbeat.DatabaseConfig{
		DBHost: "127.0.0.1",
		DBname: "plugis",
		DBuser: "plugis",
		DBpass: "plugis",
		DBPort: 3306,
	}

	servicesRepository.Start(ctx, heartbeat.NewHeartbeatSaver(dbConfig))

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("heartbeat-saver service ending")
}
