package grpool

import "context"

// GrPool --
type GrPool interface{}

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

	return gr
}

func (gp *grPool) Sync(ctx context.Context, runner Runner) error {
	return gp.interceptor(ctx, runner)
}

func (gp *grPool) Async(runner Runner) {
}
