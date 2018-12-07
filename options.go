package grpool

// Option --
type Option func(*grPool)

// WithPoolSize --
func WithPoolSize(size uint) Option {
	return func(g *grPool) {
		g.size = int(size)
	}
}

// WithUnaryInterceptor --
func WithUnaryInterceptor(interceptor Interceptor) Option {
	return func(g *grPool) {
		g.interceptor = interceptor
	}
}
