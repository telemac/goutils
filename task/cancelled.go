package task

import "context"

// IsCancelled returns true if the context is cancelled, false otherwise, without blocking
func IsCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
	return false
}
