package middlewares

import (
	"github.com/hlts2/grpool"
)

// ChainInterceptors creates a single interceptor out of a chain of many interceptors.
// For example ChainInterceptors(one, two, three) will execute one before two before three, and three.
func ChainInterceptors(interceptors ...grpool.Interceptor) grpool.Interceptor {
	n := len(interceptors)

	if n > 1 {
		lastIdx := n - 1

		return func(job grpool.Job) error {
			var (
				idx      int
				chainJob grpool.Job
			)

			chainJob = func() error {
				if idx == lastIdx {
					return job()
				}

				idx++
				return interceptors[idx](chainJob)
			}

			return interceptors[0](chainJob)
		}
	}

	if n == 1 {
		return interceptors[0]
	}

	// n == 0; Dummy interceptor
	return func(job grpool.Job) error {
		return job()
	}
}
