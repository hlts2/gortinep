package grpool

// Option --
type Option func(*grPool)

// WithPoolSize --
func WithPoolSize(size int) Option {
	return func(gp *grPool) {
		if size < 1 {
			return
		}
		gp.poolSize = size
		gp.workers = make([]*worker, 0, gp.poolSize)
	}
}

// WithUnaryInterceptor --
func WithUnaryInterceptor(interceptor Interceptor) Option {
	return func(gp *grPool) {
		gp.interceptor = interceptor
	}
}
