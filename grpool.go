package grpool

import "context"

// DefaultPoolSize --
const DefaultPoolSize = 100

// GrPool --
type GrPool interface {
	Sync(ctx context.Context, runner Runner) error
}

type grPool struct {
	size        int
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

	return gr
}

func (gp *grPool) Sync(ctx context.Context, runner Runner) error {
	return gp.interceptor(ctx, runner)
}

func (gp *grPool) Async(runner Runner) {
}
