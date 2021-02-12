package main

import (
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/shell"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"time"
)

func main() {
	// create main context
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	servicesRepository, err := natsservice.NewNatsServiceRepository("sample", "https://nats1.plugis.com", "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("sample service starting")

	// create  service
	servicesRepository.Start(ctx, &shell.ShellService{})

	cmdEvent := event.New(event.CloudEventsVersionV1)
	cmdEvent.SetType("com.plugis.shell.command")
	cmd := shell.ShellCommandParams{Command: []string{"df", "-lh", "/"}}
	cmdEvent.SetData(event.ApplicationJSON, cmd)
	cmdEvent.SetID("my-unique-id")

	resp, err := servicesRepository.Transport().Request(ctx, cmdEvent, "com.plugis.shell", time.Second*5)
	if err != nil {
		servicesRepository.Logger().WithError(err).Warn("request failed")
	}
	fmt.Printf("resp = %+v", resp)

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("tempest service ending")

}
