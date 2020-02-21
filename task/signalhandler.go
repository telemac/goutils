package task

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler allows to cancel a cancellable context on os interrupt signal reception
type SignalHandler struct {
	cancelFunc context.CancelFunc
}

// Run takes a context and a context.CancelFunc as second parameter.
// Wned the process receives an interrupt signal, the cancel function is called.
// Run is non blocking, it creates a goroutine to catch
func (sh *SignalHandler) Run(ctx context.Context, params ...interface{}) error {
	if len(params) != 1 {
		return errors.New("must have one argument")
	}
	var ok bool
	sh.cancelFunc, ok = params[0].(context.CancelFunc)
	if !ok {
		return errors.New("first argument must be of type context.CancelFunc")
	}

	signaled := make(chan os.Signal, 1)
	signal.Notify(signaled, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer sh.cancelFunc()
		select {
		case <-signaled:
		case <-ctx.Done():
		}
	}()
	return nil
}
