package heartbeat

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/natsevents"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
)

type HeartbeatSender struct {
	natsservice.NatsService
	Period        int
	RandomPeriod  int
	sentEventData *Sent
}

func NewHeartbeatSender(period int, randomPeriod int) *HeartbeatSender {
	return &HeartbeatSender{Period: period, RandomPeriod: randomPeriod}
}

func (svc *HeartbeatSender) Logger() *logrus.Entry {
	return svc.NatsService.Logger().WithField("nats-service", reflect.TypeOf(*svc).String())
}

func (svc *HeartbeatSender) SendHeartbeatEvent(ctx context.Context) error {
	t := svc.Transport()
	var err error

	// update event data field
	svc.sentEventData.Uptime = uint64(time.Since(svc.sentEventData.Started).Seconds())

	heartbeatEvent := natsevents.NewEvent("com.plugis.", "", svc.sentEventData)
	topic := heartbeatEvent.Type() + "." + svc.sentEventData.Mac
	err = t.Send(ctx, heartbeatEvent, topic)
	svc.Logger().WithFields(logrus.Fields{"event": heartbeatEvent, "topic": topic}).Trace("send event")
	if err != nil {
		svc.Logger().WithError(err).WithField("heartbeat-event", heartbeatEvent).Warn("send heartbeat cloud event")
	}
	return err
}

func (svc *HeartbeatSender) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("heartbeat sender started")
	defer log.Debug("heartbeat sender ended")

	var err error

	svc.sentEventData, err = NewSent(reflect.TypeOf(*svc).String())
	if err != nil {
		log.WithError(err).Errorf("create heartbeat.Sent event")
		return err
	}

	failureCount := 0
	for {
		err = svc.SendHeartbeatEvent(ctx)
		if err != nil {
			failureCount++
			if failureCount > 3 {
				log.WithError(err).Error("too many consecutive heartbeat failed, exit process")
				time.Sleep(time.Second * 3)
				os.Exit(1)
			}
		} else {
			failureCount = 0
		}
		waitTime := time.Second * time.Duration(svc.Period+rand.Intn(svc.RandomPeriod))
		interrupted := task.Sleep(ctx, waitTime)
		if interrupted {
			//ctx2, _ := context.WithTimeout(context.TODO(), time.Second*5)
			return svc.SendHeartbeatEvent(ctx)
		}
	}

}
