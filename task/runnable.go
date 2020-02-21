package task

import "context"

// Runnable implements a Run function
type Runnable interface {
	Run(ctx context.Context, params ...interface{}) error
}
