package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/cmd/event-sender/config"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	var params config.EventSenderConfig
	params.Parse()

	servicesRepository, err := natsservice.NewNatsServiceRepository("event-sender", params.Server, params.LogLevel)
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Infof("service starting")
	defer servicesRepository.Logger().Infof("service ended")

	eventType := params.EventType
	var obj interface{}
	err = json.Unmarshal([]byte(params.EventData), &obj)
	if err != nil {
		servicesRepository.Logger().WithError(err).Error("decode json")
	}

	topic := params.Topic

	cloudEvent := servicesRepository.Transport().NewEvent("", eventType, obj)

	if params.Request {
		ev, err := servicesRepository.Transport().Request(ctx, cloudEvent, topic, time.Second*60)
		_ = ev
		if err != nil {
			servicesRepository.Logger().WithError(err).WithField("event-type", eventType).Warn("request cloud event")
		} else {
			fmt.Printf("ev = %+v", ev)
		}
	} else {
		err = servicesRepository.Transport().Send(ctx, cloudEvent, topic)
		if err != nil {
			servicesRepository.Logger().WithError(err).WithField("event-type", eventType).Warn("send cloud event")
		}
	}

}
