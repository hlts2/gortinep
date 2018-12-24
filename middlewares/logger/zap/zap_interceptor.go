package gorpool_zap

import (
	"context"
	"time"

	"github.com/hlts2/gorpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Interceptor returns a new interceptor for zap logger.
func Interceptor(logger *zap.Logger) gorpool.Interceptor {
	return func(ctx context.Context, job gorpool.Job) error {
		logger.Info("start job")

		startT := time.Now()

		err := job(ctx)

		logger.Check(level(err), "finish job").Write(
			zap.Error(err),
			zap.Duration("time", time.Since(startT)),
		)

		return err
	}
}

func level(err error) zapcore.Level {
	if err != nil {
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}
