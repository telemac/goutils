package heartbeat

import (
	"context"
	"errors"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/natsservice"
	"reflect"
)

// HeartbeatSaver listens to com.plugis.heartbeat.Sent events
// and saves events in a database
type HeartbeatSaver struct {
	natsservice.NatsService
	DatabaseConfig DatabaseConfig
	db             Database
}

func NewHeartbeatSaver(databaseConfig DatabaseConfig) *HeartbeatSaver {
	return &HeartbeatSaver{DatabaseConfig: databaseConfig}
}

func (svc *HeartbeatSaver) Logger() *logrus.Entry {
	return svc.NatsService.Logger().WithField("nats-service", reflect.TypeOf(*svc).String())
}

func (svc *HeartbeatSaver) eventHandler(topic string, receivedEvent *event.Event, payload []byte, err error) (*event.Event, error) {
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
	case "com.plugis.heartbeat.Sent":
		var heartbeatSent Sent
		err = receivedEvent.DataAs(&heartbeatSent)
		if err != nil {
			svc.Logger().WithError(err).WithField("type", reflect.TypeOf(heartbeatSent).String()).Warn("decode event")
			return nil, err
		}

		svc.Logger().WithFields(logrus.Fields{
			"mac":      heartbeatSent.Mac,
			"hostname": heartbeatSent.Hostname,
			"ip":       heartbeatSent.InternalIP,
			"uptime":   heartbeatSent.Uptime,
		}).Info("received heartbeat")

		// TODO : save heartbeat to database
		err = svc.db.upsertHeartbeat(heartbeatSent)
		if err != nil {
			svc.Logger().WithError(err).Error("save heartbeat to database")
		}

		return nil, err

	default:
		svc.Logger().WithFields(logrus.Fields{
			"topic": topic,
			"type":  receivedEvent.Type(),
			"event": receivedEvent,
		}).Warn("unknown event type")
	}

	return nil, errors.New("unattainable code")
}

func (svc *HeartbeatSaver) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("heartbeat-saver service started")
	defer log.Debug("heartbeat-saver service ended")

	err := svc.db.Open(svc.DatabaseConfig)
	if err != nil {
		log.WithError(err).Error("connect to database")
	}

	// register eventHandler for event reception
	err = svc.Transport().RegisterHandler(svc.eventHandler, "com.plugis.heartbeat.>")
	if err != nil {
		log.WithError(err).Error("failed to register event handler")
		return err
	}

	<-ctx.Done()
	return nil
}
