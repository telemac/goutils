package ansible

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/ansibleutils"
	"github.com/telemac/goutils/logger"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/net"
	"time"
)

type AnsibleCommandParams struct {
	Command  []string `json:"command"`
	Response string   `json:"response,omitempty"`
	Error    error    `json:"error,omitempty"`
}

type AnsibleService struct {
	natsservice.NatsService
}

func (svc *AnsibleService) eventHandler(topic string, receivedEvent *event.Event, payload []byte, err error) (*event.Event, error) {
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
	case "com.plugis.ansible.playbook":
		var ansibleParams ansibleutils.AnsibleParams
		err = receivedEvent.DataAs(&ansibleParams)
		if err != nil {
			svc.Logger().WithError(err).Warn("decode AnsibleParams")
			return nil, err
		}

		playbookTimeout := time.Second * 30

		extensions := receivedEvent.Extensions()
		timeout, ok := extensions["timeout"]
		if ok {
			timeoutInt, ok := timeout.(int32)
			if ok {
				playbookTimeout = time.Duration(timeoutInt) * time.Second
			}
		}
		ctx,cancel := context.WithTimeout(context.TODO(),playbookTimeout)
		defer cancel()
		ctx = logger.WithLogger(ctx, svc.Logger())
		a := ansibleutils.New()
		result, err := a.RunPlaybooks(ctx, ansibleParams)
		if err != nil {
			svc.Logger().WithError(err).Error("run playbook")
		}

		//responseEvent := svc.Transport().NewEvent(receivedEvent.Type(), ".response", result)
		//responseEvent.SetData(cloudevents.ApplicationJSON, result)

		responseEvent := natsevents.NewEvent(receivedEvent.Type(), ".response", result)
		responseEvent.SetSource("com.plugis.ansible.playbook")

		return responseEvent, err


		// recycle received event to respond
		if result != nil {
			receivedEvent.SetExtension("response", "result[0]")
		}

		receivedEvent.SetTime(time.Now())
		receivedEvent.SetID(uuid.NewString())
		mac, _ := net.GetMACAddress()
		_ = receivedEvent.SetData(cloudevents.ApplicationJSON, "ma réponse de "+mac)

		return receivedEvent, err
	default:
		return nil, fmt.Errorf("unknown event type %s", receivedEvent.Type())
	}

	// return nil, errors.New("unattainable code")
}

func (svc AnsibleService) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()

	log.Debug("AnsibleService started")
	defer log.Debug("AnsibleService ended")

	var err error

	// register eventHandler for event reception
	mac, err := net.GetMACAddress()
	if err != nil {
		svc.Logger().WithError(err).Error("get mac accress")
	}
	topic := "com.plugis.ansible." + mac
	err = svc.Transport().RegisterHandler(svc.eventHandler, topic)
	if err != nil {
		svc.Logger().WithError(err).Error("register event handler")
		return err
	}

	<-ctx.Done()

	return nil
}
