package task

import "context"

// Runnable implements a Run function
type Runnable interface {
	Run(ctx context.Context, params ...interface{}) error
}

// IsRunnable returns true if the object implements the Runnable interface
func IsRunnable(obj interface{}) bool {
	_, ok := obj.(Runnable)
	return ok
}