package grpool

// Option --
type Option func(*grPool)

// WithGrPoolSize --
func WithGrPoolSize(size int) Option {
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
