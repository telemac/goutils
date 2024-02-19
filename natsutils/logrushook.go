package natsutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type NATSHook struct {
	config HookConfig
}

type HookConfig struct {
	NatsUtils        *NatsUtils
	Logger           *logrus.Logger
	AdditionalFields logrus.Fields
	Service          string
	Customer         string
}

func NewNatsHook(config HookConfig) (*NATSHook, error) {
	hook := &NATSHook{config: config}

	var err error

	if hook.config.NatsUtils == nil {
		return nil, errors.New("natsutils is nil")
	}
	if hook.config.NatsUtils.js == nil {
		return nil, errors.New("jetstream is nil")
	}
	if hook.config.Logger == nil {
		return nil, errors.New("Logger is nil")
	}
	if hook.config.Service == "" {
		return nil, errors.New("Service is empty")
	}
	if hook.config.Customer == "" {
		return nil, errors.New("Customer is empty")
	}

	// replace . by _ in customer name
	hook.config.Customer = strings.Replace(hook.config.Customer, ".", "_", -1)
	// replace . by _ in service name
	hook.config.Service = strings.Replace(hook.config.Service, ".", "_", -1)

	s, err := hook.config.NatsUtils.js.CreateStream(context.TODO(), jetstream.StreamConfig{
		Name:        "LOGS",
		Description: "micro-service logs",
		Subjects:    []string{"log.>"},
		MaxAge:      time.Hour * 24 * 30,
	})
	_ = s
	if err != nil {
		hook.config.Logger.WithError(err).WithField("function", "NewNATSHook").Warn("create LOGS stream")
		//return nil, err
	}

	hook.config.Logger.SetFormatter(&logrus.JSONFormatter{})
	//log.Logger.SetOutput(os.Stderr)
	hook.config.Logger.AddHook(hook)

	return hook, nil
}

func (hook *NATSHook) Fire(entry *logrus.Entry) error {
	var err error
	var msg string

	var topic string = fmt.Sprintf("log.%s.%s.%s", entry.Level.String(), hook.config.Service, hook.config.Customer)

	if hook.config.AdditionalFields == nil {
		msg, err = entry.String()
	} else {
		newEntry := *entry.Dup().WithFields(hook.config.AdditionalFields)
		newEntry.Level = entry.Level
		newEntry.Message = entry.Message
		msg, err = newEntry.String()
	}

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	hook.config.NatsUtils.PublishOnStream(ctx, "LOGS", topic, []byte(msg))
	hook.config.NatsUtils.Nats().Flush()

	return nil
}

func (hook *NATSHook) Levels() []logrus.Level {
	return logrus.AllLevels
	//return []logrus.Level{
	//	logrus.PanicLevel,
	//	logrus.FatalLevel,
	//	logrus.ErrorLevel,
	//	logrus.WarnLevel,
	//	logrus.InfoLevel,
	//	logrus.DebugLevel,
	//	logrus.TraceLevel,
	//}
}
