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
	gp := &grPool{
		size:       DefaultPoolSize,
		runnableCh: make(chan runnable),
	}

	for _, opt := range opts {
		opt(gp)
	}

	for i := 0; i < gp.size; i++ {
		go gp.async()
	}

	return gp
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

func (gp *grPool) async() {
	for r := range gp.runnableCh {
		if gp.interceptor == nil {
			r.runner(r.ctx)
		} else {
			gp.interceptor(r.ctx, r.runner)
		}
	}
}
