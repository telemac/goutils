package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/tcpgetcommand"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
)

func main() {
	serviceRepositoryName := "tcp-get-command"
	config, err := natsservice.LoadConfig("servers.yml")
	if err != nil {
		logrus.WithError(err).Fatal("open config file")
	}

	// create main context
	ctx, cancel := task.NewCancellableContext(time.Second * 5)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository(serviceRepositoryName, config.Servers[0].Url, config.CommandLineParams.Log)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"servers": config.Servers[0].Url,
			"error":   err,
		}).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info(serviceRepositoryName + " service starting")

	servicesRepository.Start(ctx, tcpgetcommand.NewTcpGetCommandService(tcpgetcommand.TcpGetCommandConfig{
		ListenPort:    10,
		ListenAddress: "0.0.0.0",
	}))

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info(serviceRepositoryName + " service ending")
}
