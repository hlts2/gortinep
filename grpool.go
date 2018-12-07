package grpool

// GrPool --
type GrPool interface{}

type grPool struct {
	size        int
	interceptor Interceptor
}

// NewGrPool --
func NewGrPool(opts ...Option) GrPool {
	gr := new(grPool)

	for _, opt := range opts {
		opt(gr)
	}

	return gr
}

func (gp *grPool) Sync(runable Runnable) error {
	return nil
}

func (gp *grPool) Async(runable Runnable) {
}
