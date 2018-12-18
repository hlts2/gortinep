package grpool

// Option configures grpool.
type Option func(*grPool)

// WithPoolSize returns an option that sets the pool size.
func WithPoolSize(size int) Option {
	return func(gp *grPool) {
		if size < 1 {
			return
		}
		gp.poolSize = size
		gp.workers = make([]*worker, 0, gp.poolSize)
	}
}

// WithError returns an option that sets channel for job error processed by goroutine worker.
// The result of each goroutine is sent to this channe.
func WithError(errCh chan error) Option {
	return func(gp *grPool) {
		gp.errCh = errCh
	}
}

// WithInterceptor returns an option that sets the Interceptor implementation.
func WithInterceptor(interceptor Interceptor) Option {
	return func(gp *grPool) {
		gp.interceptor = interceptor
	}
}
