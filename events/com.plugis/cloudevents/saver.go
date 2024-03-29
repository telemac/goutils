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
	postgresConfig natsservice.PostgresConfig
	mysqlConfig    natsservice.MysqlConfig
	postgresDb     CloudEventsDatabase
	mysqlDb        VariablesDatabase
}

func NewCloudEventSaver(postgresConfig natsservice.PostgresConfig, mysqlConfig natsservice.MysqlConfig) *CloudEventSaver {
	return &CloudEventSaver{
		postgresConfig: postgresConfig,
		mysqlConfig:    mysqlConfig,
	}
}

func (svc *CloudEventSaver) Logger() *logrus.Entry {
	return svc.NatsService.Logger().WithField("nats-service", reflect.TypeOf(*svc).String())
}

func (svc *CloudEventSaver) eventHandler(topic string, receivedEvent *event.Event, payload []byte, CEDecodeErr error) (*event.Event, error) {
	// save cloudevent to postgres cloudevents table

	logger := svc.Logger().WithFields(logrus.Fields{
		"topic":   topic,
		"event":   receivedEvent,
		"payload": string(payload),
	})

	if receivedEvent == nil {
		logger.WithError(CEDecodeErr).Warn("unable to decode event")
		return nil, nil
	}

	err := svc.postgresDb.InsertEvent(topic, receivedEvent, payload, CEDecodeErr)
	logger.WithError(err).Trace("received event")
	if err != nil {
		logger.WithError(err).Error("log cloudevent to postgres cloudevents table")
	}
	/*
		// if event is variable set, save value to mysql variable table
		if receivedEvent != nil && receivedEvent.Type() == "com.plugis.variable.set" {
			logger.Info("variable set")
			var variables variable.Variables
			err = receivedEvent.DataAs(&variables)
			if err != nil {
				logger.WithError(err).Warn("decode variable set data")
				return nil, nil
			}
			err = svc.mysqlDb.upsertVariables(receivedEvent.ID(), variables)
			if err != nil {
				logger.WithError(err).Warn("upsert variables")
				return nil, nil
			}
		}
	*/
	// don't send a response to any request
	return nil, nil
}

func (svc *CloudEventSaver) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("service started")
	defer log.Debug("service ended")

	// open cloudevents database (postgreSql)
	err := svc.postgresDb.Open(svc.postgresConfig)
	if err != nil {
		log.WithError(err).Error("connect to PostgreSQL database")
		return err
	}

	/*
		// open variables database (mysql)
		err = svc.mysqlDb.Open(svc.mysqlConfig)
		if err != nil {
			log.WithError(err).Error("connect to MySQL database")
			return err
		}
	*/

	// register eventHandler for event reception
	err = svc.Transport().RegisterHandler(svc.eventHandler, ">")
	if err != nil {
		log.WithError(err).Error("failed to register event handler")
		return err
	}

	//for !task.IsCancelled(ctx) {
	//	err = svc.postgresDb.Cleanheartbeats()
	//	if err != nil {
	//		log.WithError(err).Error("clean heartbeat cloudevents")
	//	}
	//	task.Sleep(ctx, time.Minute)
	//}

	<-ctx.Done()
	return nil
}
