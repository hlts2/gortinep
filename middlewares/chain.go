package middlewares

import (
	"github.com/hlts2/grpool"
)

// ChainUnaryInterceptors --
func ChainUnaryInterceptors(interceptors ...grpool.Interceptor) grpool.Interceptor {
	return func(runner grpool.Runner) error {
		var (
			idx         int
			chainRunner grpool.Runner
		)

		chainRunner = func() error {
			if idx == len(interceptors) {
				return runner()
			}
			idx++
			return interceptors[idx](chainRunner)
		}

		return interceptors[0](chainRunner)
	}
}
