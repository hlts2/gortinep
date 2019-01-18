package gortinep_log

import (
	"context"
	"time"

	"github.com/hlts2/gortinep"
)

// Interceptor returns a new interceptor for log.
func Interceptor(ops ...Option) gortinep.Interceptor {
	var o = evaluateOption(ops...)

	return func(ctx context.Context, job gortinep.Job) error {
		startTime := time.Now()

		o.logger.Print("start job.")

		err := job(ctx)

		o.logger.Printf("finish job. err: %v, time: %v", err, time.Since(startTime))

		return err
	}
}
