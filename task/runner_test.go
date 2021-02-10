package task

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type testService struct {
}

func (svc *testService) Run(ctx context.Context, params ...interface{}) error {
	count := 0
	maxCount, ok := params[0].(int)
	if !ok {
		return fmt.Errorf("maxCount not an int %s", reflect.TypeOf(params))
	}
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			count++
			fmt.Printf("tick %d/%d\n", count, maxCount)
			if count == maxCount {
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

func TestRunnerRepository_Run(t *testing.T) {
	//assert := assert.New(t)
	ctx, cancel := NewCancellableContext(time.Second * 10)
	defer cancel()

	services := NewRunnerRepository()

	ts := new(testService)

	t1 := services.Start(ctx, ts, 3)
	time.Sleep(time.Millisecond * 100)
	t2 := services.Start(ctx, ts, 5)

	services.WaitUntilAllDone()

	fmt.Printf("errChan1 = %v\n", <-t1.err)
	fmt.Printf("errChan2 = %v\n", <-t2.err)

}
