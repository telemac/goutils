package task

import (
	"context"
	"time"
)

// Sleep is an interruptible time.Sleep version
// It returns true if it was interrupted by context cancellation, false otherwise
func Sleep(ctx context.Context, duration time.Duration) (wasInterrupted bool) {
	select {
	case <-ctx.Done():
		wasInterrupted = true
	case <-time.After(duration):
		wasInterrupted = false
	}
	return
}
