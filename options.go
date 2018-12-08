package grpool

// Option --
type Option func(*grPool)

// WithPoolSize --
func WithPoolSize(size int) Option {
	return func(g *grPool) {
		if size < 1 {
			return
		}
		g.poolSize = size
		g.workers = make([]*worker, 0, g.poolSize)
	}
}

// WithUnaryInterceptor --
func WithUnaryInterceptor(interceptor Interceptor) Option {
	return func(g *grPool) {
		g.interceptor = interceptor
	}
}
