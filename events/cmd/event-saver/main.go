package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/cloudevents"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

var Config = struct {
	APPName string `default:"event-saver"`

	DB struct {
		DBType   string `default:"postgre"`
		Host     string `default:"127.0.0.1"`
		Database string `default:"plugis"`
		User     string `default:"plugis"`
		Password string `required:"true" env:"DBPassword" default:"plugis"`
		Port     uint   `default:"5432"`
	}

	Servers []struct {
		url string `default:"nats://nats1.plugis.com:443"`
		sni bool   `default:"false"`
	}
}{}

func main() {
	err := configor.Load(&Config, "./event-saver.yml")
	if err != nil {
		logrus.WithError(err).Warn("load configuration file")
	}

	fmt.Printf("config: %#v", Config)

	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("event-saver", "nats://nats1.plugis.com:443", "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("service starting")

	// start event-saver service
	servicesRepository.Start(ctx, &cloudevents.CloudEventSaver{})

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("service ending")
}
