package main

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/browser"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	config, err := natsservice.LoadConfig("./remote-access-dashboard.yml")
	if err != nil {
		logrus.WithError(err).Fatal("open config file")
	}

	// create main context
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("remote-access-dashboard", "https://nats1.plugis.com", "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("remote-access-dashboard service starting")

	servicesRepository.Start(ctx, heartbeat.NewHeartbeatWebInterface(config.Mysql))

	servicesRepository.Start(ctx, &browser.BrowserService{})

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("remote-access-dashboard service ending")
}
