package main

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	config, err := natsservice.LoadConfig("servers.yml", "mysql.yml")
	if err != nil {
		logrus.WithError(err).Fatal("open config file")
	}

	ctx, cancel := task.NewCancellableContext(time.Second * 10)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("heartbeat-saver", config.Servers[0].Url, config.CommandLineParams.Log)
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("heartbeat-saver service starting")

	// start heartbeat saver
	servicesRepository.Start(ctx, heartbeat.NewHeartbeatSaver(config.Mysql))

	// start heartbeat web interface
	//servicesRepository.Start(ctx, heartbeat.NewHeartbeatWebInterface(config.Mysql))

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("heartbeat-saver service ending")
}
