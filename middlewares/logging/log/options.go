package grpool_log

import (
	"log"
	"os"
)

type option struct {
	logger Logger
}

// Option configures grpool_log
type Option func(*option)

func evaluateOption(ops ...Option) *option {
	eo := &option{
		logger: &basicLogger{
			l: log.New(os.Stdout, "", log.LstdFlags),
		},
	}

	for _, op := range ops {
		op(eo)
	}

	return eo
}

// WithLogger returns an option that sets Logger implementation.
func WithLogger(logger Logger) Option {
	return func(o *option) {
		if logger == nil {
			return
		}
		o.logger = logger
	}
}
