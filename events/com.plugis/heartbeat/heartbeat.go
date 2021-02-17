package heartbeat

import (
	"context"
	"time"

	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
)

type HeartbeatService struct {
	natsservice.NatsService
	Period       int
	RandomPeriod int
	sentEvent    *Sent
}

func (svc *HeartbeatService) SendHeartbeatEvent(ctx context.Context) error {
	t := svc.Transport()
	var err error

	svc.sentEvent, err = NewSent()
	if err != nil {
		svc.Logger().WithError(err).Errorf("create heartbeat.Sent event")
	}

	heartbeatEvent := t.NewEvent("com.plugis.", "", svc.sentEvent)
	err = t.Send(ctx, heartbeatEvent, heartbeatEvent.Type()+"."+svc.sentEvent.Mac)
	if err != nil {
		svc.Logger().WithError(err).WithField("heartbeat-event", heartbeatEvent).Warn("send heartbeat cloud event")
	}
	return err
}

func (svc *HeartbeatService) Run(ctx context.Context, params ...interface{}) error {
	svc.Logger().Debug("heartbeat service started")
	defer svc.Logger().Debug("heartbeat service ended")

	for {
		_ = svc.SendHeartbeatEvent(ctx)
		interrupted := task.Sleep(ctx, time.Second*time.Duration(svc.Period))
		if interrupted {
			//ctx2, _ := context.WithTimeout(context.TODO(), time.Second*5)
			return svc.SendHeartbeatEvent(ctx)
		}
	}

}
