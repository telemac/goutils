package main

import (
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/cloudevents"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
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
