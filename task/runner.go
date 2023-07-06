package task

import (
	"context"
	"sync"
)

type Runner interface {
	Start(ctx context.Context, runnable Runnable, params ...interface{}) *Task
	WaitUntilAllDone()
}

type Task struct {
	runnable Runnable
	params   interface{}
	err      chan error
}

type RunnerRepository struct {
	wg    sync.WaitGroup
	tasks []*Task
	mutex sync.RWMutex
}

func NewRunnerRepository() *RunnerRepository {
	return &RunnerRepository{}
}

func (repo *RunnerRepository) Start(ctx context.Context, runnable Runnable, params ...interface{}) *Task {
	// check if runnable is an instance of Runnable
	_, ok := runnable.(Runnable)
	if !ok {
		panic("runnable must be not nil and a task.Runnable")
	}
	t := new(Task)
	t.runnable = runnable
	t.params = params
	t.err = make(chan error, 1)
	repo.mutex.Lock()
	repo.tasks = append(repo.tasks, t)
	repo.mutex.Unlock()
	repo.wg.Add(1)
	go func(ctx context.Context, task *Task) {
		task.err <- runnable.Run(ctx, params...)
		repo.wg.Done()
	}(ctx, t)
	return t
}

func (repo *RunnerRepository) WaitUntilAllDone() {
	repo.wg.Wait()
}
