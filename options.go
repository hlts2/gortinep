package gortinep

// Option configures gortinep.
type Option func(*grPool)

// WithPoolSize returns an option that sets the pool size.
func WithPoolSize(size int) Option {
	return func(gp *grPool) {
		if size < 1 {
			return
		}
		gp.poolSize = size
		gp.workers = make([]*worker, gp.poolSize)
	}
}

// WithJobSize returns an option that sets the job size.
func WithJobSize(size int) Option {
	return func(gp *grPool) {
		if size < 1 {
			return
		}
		gp.jobSize = size
		gp.jobCh = make(chan Job, size)
	}
}

// WithErrorChannel returns an option that sets channel for job error processed by goroutine worker.
// The result of each goroutine is sent to this channe.
func WithErrorChannel(ch chan error) Option {
	return func(gp *grPool) {
		if ch == nil {
			return
		}
		gp.wrapperrCh = &wrapperrCh{
			ch:     ch,
			closed: false,
		}
	}
}

// WithInterceptor returns an option that sets the Interceptor implementation.
func WithInterceptor(interceptor Interceptor) Option {
	return func(gp *grPool) {
		gp.interceptor = interceptor
	}
}
