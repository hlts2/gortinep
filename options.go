package grpool

// Option --
type Option func(*grPool)

// WithPoolSize --
func WithPoolSize(size int) Option {
	return func(g *grPool) {
		g.size = size
	}
}

// WithUnaryInterceptor --
func WithUnaryInterceptor(interceptor Interceptor) Option {
	return func(g *grPool) {
		g.interceptor = interceptor
	}
}
