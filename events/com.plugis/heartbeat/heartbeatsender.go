package heartbeat

import (
	"context"
	"math/rand"
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

func (svc *HeartbeatSender) SendHeartbeatEvent(ctx context.Context) error {
	t := svc.Transport()
	var err error

	// update event data field
	svc.sentEventData.Uptime = uint64(time.Since(svc.sentEventData.Started).Seconds())

	heartbeatEvent := t.NewEvent("com.plugis.", "", svc.sentEventData)
	err = t.Send(ctx, heartbeatEvent, heartbeatEvent.Type()+"."+svc.sentEventData.Mac)
	if err != nil {
		svc.Logger().WithError(err).WithField("heartbeat-event", heartbeatEvent).Warn("send heartbeat cloud event")
	}
	return err
}

func (svc *HeartbeatSender) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger().WithField("type", reflect.TypeOf(svc).String())
	log.Debug("heartbeat service started")
	defer log.Debug("heartbeat service ended")

	var err error

	svc.sentEventData, err = NewSent()
	if err != nil {
		log.WithError(err).Errorf("create heartbeat.Sent event")
		return err
	}

	for {
		_ = svc.SendHeartbeatEvent(ctx)

		waitTime := time.Second * time.Duration(svc.Period+rand.Intn(svc.RandomPeriod))
		interrupted := task.Sleep(ctx, waitTime)
		if interrupted {
			//ctx2, _ := context.WithTimeout(context.TODO(), time.Second*5)
			return svc.SendHeartbeatEvent(ctx)
		}
	}

}
