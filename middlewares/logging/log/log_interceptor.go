package grpool_log

import (
	"context"
	"sync"

	"github.com/hlts2/grpool"
)

// Interceptor returns a new interceptor for log.
func Interceptor(ops ...Option) grpool.Interceptor {
	var (
		o  = evaluateOption(ops...)
		mu = new(sync.Mutex)
	)
	return func(ctx context.Context, job grpool.Job) error {
		mu.Lock()
		defer mu.Unlock()

		o.logger.Printf("start job.")

		_ = job(ctx)

		o.logger.Printf("finish job.")

		return nil
	}
}
