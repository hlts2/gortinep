package grpool_log

import (
	"sync"

	"github.com/hlts2/grpool"
)

// Interceptor returns a new interceptor for log.
func Interceptor(ops ...Option) grpool.Interceptor {
	var (
		o  = evaluateOption(ops...)
		mu = new(sync.Mutex)
	)
	return func(job grpool.Job) error {
		mu.Lock()
		defer mu.Unlock()

		o.logger.Printf("start job.")

		_ = job()

		o.logger.Printf("finish job.")

		return nil
	}
}
