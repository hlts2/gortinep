package middlewares

import (
	"context"

	"github.com/hlts2/gorpool"
)

// ChainInterceptors creates a single interceptor out of a chain of many interceptors.
// For example ChainInterceptors(one, two, three) will execute one before two before three, and three.
func ChainInterceptors(interceptors ...gorpool.Interceptor) gorpool.Interceptor {
	n := len(interceptors)

	if n > 1 {
		lastIdx := n - 1

		return func(ctx context.Context, job gorpool.Job) error {
			var (
				idx      int
				chainJob gorpool.Job
			)

			chainJob = func(ctx context.Context) error {
				if idx == lastIdx {
					return job(ctx)
				}

				idx++
				return interceptors[idx](ctx, chainJob)
			}

			return interceptors[0](ctx, chainJob)
		}
	}

	if n == 1 {
		return interceptors[0]
	}

	// n == 0; Dummy interceptor
	return func(ctx context.Context, job gorpool.Job) error {
		return job(ctx)
	}
}
