package middlewares

import (
	"context"

	"github.com/hlts2/grpool"
)

// ChainUnaryInterceptors --
func ChainUnaryInterceptors(interceptors ...grpool.Interceptor) grpool.Interceptor {
	return func(ctx context.Context, handler grpool.Handler) {
		var (
			idx          int
			chainHandler grpool.Handler
		)

		chainHandler = func(ctx context.Context) {
			if idx == len(interceptors) {
				handler(ctx)
				return
			}
			idx++
			interceptors[idx](ctx, chainHandler)
		}

		interceptors[0](ctx, chainHandler)
	}
}
