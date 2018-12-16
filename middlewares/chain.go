package middlewares

import (
	"github.com/hlts2/grpool"
)

// ChainUnaryInterceptors creates a single interceptor out of a chain of many interceptors.
// For example ChainUnaryInterceptors(one, two, three) will execute one before two before three, and three.
func ChainUnaryInterceptors(interceptors ...grpool.Interceptor) grpool.Interceptor {
	return func(job grpool.Job) error {
		var (
			idx      int
			chainJob grpool.Job
		)

		chainJob = func() error {
			if idx == len(interceptors) {
				return job()
			}
			idx++
			return interceptors[idx](chainJob)
		}

		return interceptors[0](chainJob)
	}
}
