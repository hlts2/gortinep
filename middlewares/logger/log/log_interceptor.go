package gorpool_log

import (
	"context"
	"sync"
	"time"

	"github.com/hlts2/gorpool"
)

// Interceptor returns a new interceptor for log.
func Interceptor(ops ...Option) gorpool.Interceptor {
	var (
		o  = evaluateOption(ops...)
		mu = new(sync.Mutex)
	)
	return func(ctx context.Context, job gorpool.Job) error {
		mu.Lock()
		defer mu.Unlock()

		startTime := time.Now()

		o.logger.Print("start job.")

		err := job(ctx)

		o.logger.Printf("finish job. err: %v, time: %v", err, time.Since(startTime))

		return err
	}
}
