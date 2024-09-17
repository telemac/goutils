package task

import "context"

// Runnable implements a Run function
type Runnable interface {
	Run(ctx context.Context, params ...interface{}) error
}

/*
// IsRunnable returns true if the object implements the Runnable interface
func IsRunnable(obj Runnable) bool {
	if obj == nil {
		return false
	}
	runnambe, ok := (obj).(Runnable)
	if ok {
		_, ok = runnambe.(interface {
			Run(ctx context.Context, params ...interface{}) error
		})
	}
	return ok
}
*/
