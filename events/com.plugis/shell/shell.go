package shell

import (
	"context"
	"errors"
	"fmt"
	"github.com/telemac/goutils/net"
	"os/exec"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/natsservice"
)

type ShellCommandParams struct {
	Command  []string `json:"command"`
	Response string   `json:"response,omitempty"`
	Error    error    `json:"error,omitempty"`
}

type ShellService struct {
	natsservice.NatsService
}

func (svc *ShellService) eventHandler(topic string, receivedEvent *event.Event, payload []byte, err error) (*event.Event, error) {
	// check if no error on cloud event formatting
	if err != nil {
		svc.Logger().WithFields(logrus.Fields{
			"topic":   topic,
			"event":   receivedEvent,
			"payload": string(payload),
			"error":   err,
		}).Error("receive cloud event")
		return nil, err
	}

	switch receivedEvent.Type() {
	case "com.plugis.shell.command":
		var params ShellCommandParams
		err = receivedEvent.DataAs(&params)
		if err != nil {
			svc.Logger().WithError(err).Warn("decode ShellCommandParams")
			return nil, err
		}

		if len(params.Command) < 1 {
			return nil, errors.New("command neets at least one parameter")
		}
		var cmd *exec.Cmd
		if len(params.Command) == 1 {
			cmd = exec.Command(params.Command[0])
		} else if len(params.Command) > 1 {
			cmd = exec.Command(params.Command[0], params.Command[1:]...)
		}

		//params.Command = append([]string{"sh", "-c"}, params.Command...)

		out, err := cmd.CombinedOutput()
		params.Response, params.Error = string(out), err
		if params.Error != nil {
			svc.Logger().WithError(err).WithField("command", params.Command).Warn("run command")
			return nil, err
		}

		fmt.Println(string(params.Response))
		receivedEvent.SetData(event.ApplicationJSON, params)
		return receivedEvent, err

	default:
		svc.Logger().WithFields(logrus.Fields{
			"topic": topic,
			"type":  receivedEvent.Type(),
			"event": receivedEvent,
		}).Warn("unknown event type")
	}

	return nil, errors.New("unattainable code")
}

func (svc ShellService) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()

	log.Debug("ShellService started")
	defer log.Debug("ShellService ended")

	var err error

	// register eventHandler for event reception
	mac, err := net.GetMACAddress()
	if err != nil {
		svc.Logger().WithError(err).Error("get mad accress")
	}
	topic := "com.plugis.shell." + mac
	err = svc.Transport().RegisterHandler(svc.eventHandler, topic)
	if err != nil {
		svc.Logger().WithError(err).Error("register event handler")
		return err
	}

	<-ctx.Done()

	return nil
}
