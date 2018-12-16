package grpool_zap

import (
	"time"

	"github.com/hlts2/grpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Interceptor returns a new interceptor for zap logger.
func Interceptor(logger *zap.Logger) grpool.Interceptor {
	return func(job grpool.Job) error {
		logger.Info("start job")

		startT := time.Now()

		err := job()

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
