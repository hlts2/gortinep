package grpool_log

import (
	"log"
)

// Logger is the interface for log output
type Logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// LoggerFunc is the function for logger.
type LogFunc func(format string, args ...interface{})

type basicLogger struct {
	l *log.Logger
}

func (b *basicLogger) Printf(format string, args ...interface{}) {
	b.l.Printf(format, args)
}

func (b *basicLogger) Errorf(format string, args ...interface{}) {
	b.l.Printf(format, args)
}
