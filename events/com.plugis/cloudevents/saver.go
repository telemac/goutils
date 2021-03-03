package cloudevents

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"reflect"

	"github.com/telemac/goutils/natsservice"
)

// CloudEventSaver listens to all cloud events
// and saves them in a postgresql database
type CloudEventSaver struct {
	natsservice.NatsService
	db Database
}

func (svc *CloudEventSaver) Logger() *logrus.Entry {
	return svc.NatsService.Logger().WithField("nats-service", reflect.TypeOf(*svc).String())
}

func (svc *CloudEventSaver) eventHandler(topic string, receivedEvent *event.Event, payload []byte, err error) (*event.Event, error) {
	// save cloudevent to database
	err = svc.db.InsertEvent(topic, receivedEvent, payload, err)
	if err != nil {
		svc.Logger().WithFields(logrus.Fields{
			"topic":   topic,
			"event":   receivedEvent,
			"payload": string(payload),
			"error":   err,
		}).Error("log cloudevent to database")
	}
	return nil, nil
}

func (svc *CloudEventSaver) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("service started")
	defer log.Debug("service ended")

	// premare database
	dbConfig := DatabaseConfig{
		DBHost: "localhost",
		DBname: "plugis",
		DBuser: "plugis",
		DBpass: "plugis",
		DBPort: 26257,
	}

	err := svc.db.Open(dbConfig)
	if err != nil {
		log.Error("connect to database")
		return err
	}

	// register eventHandler for event reception
	err = svc.Transport().RegisterHandler(svc.eventHandler, ">")
	if err != nil {
		log.WithError(err).Error("failed to register event handler")
		return err
	}

	<-ctx.Done()
	return nil
}
