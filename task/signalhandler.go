package task

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// SignalHandler allows to cancel a cancellable context on os interrupt signal reception
type SignalHandler struct {
	cancelFunc context.CancelFunc
}

// Run takes a context and a context.CancelFunc as second parameter.
// if the third parameter is a time.Duration, the program forces exit after that timeout
// Wned the process receives an interrupt signal, the cancel function is called.
// Run is non blocking, it creates a goroutine to catch interrupt signal
func (sh *SignalHandler) Run(ctx context.Context, params ...interface{}) error {
	// get first parameter, the cancel fund
	if len(params) < 1 {
		return errors.New("must have one argument")
	}
	var ok bool
	sh.cancelFunc, ok = params[0].(context.CancelFunc)
	if !ok {
		return errors.New("first argument must be of type context.CancelFunc")
	}

	// get second parameter, the delay before forced exit
	var delayBeforeForcedExit time.Duration
	if len(params) >= 2 {
		delay, ok := params[1].(time.Duration)
		if ok {
			delayBeforeForcedExit = delay
		}
	}

	signaled := make(chan os.Signal, 1)
	signal.Notify(signaled, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer sh.cancelFunc()
		select {
		case <-signaled:
			if delayBeforeForcedExit > 0 {
				go func() {
					time.Sleep(delayBeforeForcedExit)
					fmt.Printf("\nForce exit after %s\n", delayBeforeForcedExit.String())
					os.Exit(1)
				}()
			}
		case <-ctx.Done():
		}
	}()
	return nil
}
