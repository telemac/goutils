package natsservice

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/goutils/task"
	"reflect"
	"testing"
	"time"
)

type NatsServiceSample struct {
	NatsService
}

func (svc *NatsServiceSample) Run(ctx context.Context, params ...interface{}) error {
	count := 0
	maxCount, ok := params[0].(int)
	if !ok {
		return errors.New("maxCount not an int")
	}
	tick := time.NewTicker(time.Second)
	type EventData struct {
		Count int
		Total int
	}
	for {
		select {
		case <-tick.C:
			count++
			fmt.Printf("tick %d/%d\n", count, maxCount)
			data := EventData{count, maxCount}
			ce := cloudevents.NewEvent()
			ce.SetType("com.plugis.sample.ce." + reflect.TypeOf(data).String())
			ce.SetID(uuid.New().String())
			ce.SetTime(time.Now())
			ce.SetData(cloudevents.ApplicationJSON, data)
			err := svc.transport.Send(ctx, ce, "com.plugis.test."+reflect.TypeOf(svc).String())
			if err != nil {
				svc.logger.WithError(err).Error("transport.Send")
			}
			if count == maxCount {
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

func TestNewNatsServiceRepository(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	natsSvcRepo, err := NewNatsServiceRepository("nats test service", "https://nats1.plugis.com", "trace")
	assert.NoError(err)

	svcSample := new(NatsServiceSample)

	natsSvcRepo.Start(ctx, svcSample, 3)
	natsSvcRepo.Start(ctx, svcSample, 5)

	natsSvcRepo.WaitUntilAllDone()

	err = natsSvcRepo.Close()
	assert.NoError(err)
}
