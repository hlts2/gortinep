package grpool

import (
	"context"
)

// ChainInterceptors --
func ChainInterceptors(interceptors ...Interceptor) Interceptor {
	return func(ctx context.Context, runner Runner) {
		var (
			idx         int
			chainRunner Runner
		)

		chainRunner = func(ctx context.Context) error {
			if idx == len(interceptors) {
				return runner(ctx)
			}

			idx++
			interceptors[idx](ctx, chainRunner)
			return nil
		}

		interceptors[0](ctx, chainRunner)
	}
}
