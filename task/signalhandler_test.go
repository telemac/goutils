package task

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

// TODO: more a example then a test
func TestPSignalHandler(t *testing.T) {
	assert := assert.New(t)

	// Cancelling contect
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create errgroup to run parallel processing
	g, gctx := errgroup.WithContext(ctx)
	_ = gctx

	// Handle interrupt signals
	var signalHandler SignalHandler
	err := signalHandler.Run(ctx, cancel, time.Second*10)
	assert.NoError(err, "signalHandler.Run")

	for i := 0; i < 10; i++ {
		delaySecond := i
		g.Go(func() error {
			interrupted := Sleep(ctx, time.Millisecond*time.Duration(delaySecond))
			Sleep(ctx, time.Second)
			if interrupted {
				return fmt.Errorf("Sleep %d interrupted", delaySecond)
			}
			fmt.Printf("Sleep %d done (%d)\n", delaySecond, runtime.NumGoroutine())
			return nil
		})
	}

	log.Println("*** All goroutines scheduled")

	err = g.Wait()
	if err != nil {
		log.WithError(err).Error("waitgroup error")
	}
	cancel()
	// Wait until interrupt signal received or cancel called
	<-ctx.Done()
}
