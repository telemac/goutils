package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/cloudevents"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	config, err := natsservice.LoadConfig("./event-saver.yml")
	//err := configor.Load(&Config, "./event-saver.yml")
	if err != nil {
		logrus.WithError(err).Warn("load configuration file")
	}

	fmt.Printf("config: %#v", config)

	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("event-saver", config.Servers[0].Url, "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("service starting")

	// start event-saver service
	servicesRepository.Start(ctx, cloudevents.NewCloudEventSaver(config.Postgres))

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("service ending")
}
