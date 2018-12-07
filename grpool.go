package grpool

import (
	"context"
)

// DefaultPoolSize --
const DefaultPoolSize = 100

// GrPool --
type GrPool interface {
	Sync(ctx context.Context, runner Runner) error
	Async(ctx context.Context, runner Runner)
}

type grPool struct {
	size   int
	taskCh chan struct {
		runner Runner
		ctx    context.Context
	}
	interceptor Interceptor
}

// New --
func New(opts ...Option) GrPool {
	gr := new(grPool)

	for _, opt := range opts {
		opt(gr)
	}

	if gr.size < 1 {
		gr.size = DefaultPoolSize
	}

	for i := 0; i < gr.size; i++ {
		go func() {
			async(gr.taskCh)
		}()
	}

	return gr
}

func (gp *grPool) Sync(ctx context.Context, runner Runner) error {
	if gp.interceptor == nil {
		return runner(ctx)
	}

	return gp.interceptor(ctx, runner)
}

func (gp *grPool) Async(ctx context.Context, runner Runner) {
	gp.taskCh <- struct {
		runner Runner
		ctx    context.Context
	}{
		runner: runner,
		ctx:    ctx,
	}
}

func async(taskCh chan struct {
	runner Runner
	ctx    context.Context
}) {
	for task := range taskCh {
		task.runner(task.ctx)
	}
}
