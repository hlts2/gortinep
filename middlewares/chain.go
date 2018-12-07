package middlewares

import (
	"context"

	"github.com/hlts2/grpool"
)

// ChainUnaryInterceptors --
func ChainUnaryInterceptors(interceptors ...grpool.Interceptor) grpool.Interceptor {
	return func(ctx context.Context, runner grpool.Runner) error {
		var (
			idx         int
			chainRunner grpool.Runner
		)

		chainRunner = func(ctx context.Context) error {
			if idx == len(interceptors) {
				return runner(ctx)
			}
			idx++
			return interceptors[idx](ctx, chainRunner)
		}

		return interceptors[0](ctx, chainRunner)
	}
}
