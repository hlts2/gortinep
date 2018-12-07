package middlewares

import (
	"context"

	"github.com/hlts2/grpool"
)

// ChainUnaryInterceptors --
func ChainUnaryInterceptors(interceptors ...grpool.Interceptor) grpool.Interceptor {
	return func(ctx context.Context, handler grpool.Runner) error {
		var (
			idx          int
			chainHandler grpool.Runner
		)

		chainHandler = func(ctx context.Context) error {
			if idx == len(interceptors) {
				return handler(ctx)
			}
			idx++
			return interceptors[idx](ctx, chainHandler)
		}

		return interceptors[0](ctx, chainHandler)
	}
}
