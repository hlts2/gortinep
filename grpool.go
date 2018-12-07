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

// Runnable --
type runnable struct {
	runner Runner
	ctx    context.Context
}

type grPool struct {
	size        int
	runnableCh  chan runnable
	interceptor Interceptor
}

// New --
func New(opts ...Option) GrPool {
	gr := &grPool{
		size:       DefaultPoolSize,
		runnableCh: make(chan runnable),
	}

	for _, opt := range opts {
		opt(gr)
	}

	for i := 0; i < gr.size; i++ {
		go async(gr.runnableCh)
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
	gp.runnableCh <- runnable{
		runner: runner,
		ctx:    ctx,
	}
}

func async(runnableCh chan runnable) {
	for r := range runnableCh {
		r.runner(r.ctx)
	}
}
