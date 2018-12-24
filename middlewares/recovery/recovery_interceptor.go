package gorpool_recovery

import (
	"context"

	"github.com/hlts2/gorpool"
)

// RecoveryHandlerFunc is a function that recovers from the panic `p`.
type RecoveryHandlerFunc func(p interface{})

// Interceptor returns a new interceptor for panic recovery.
func Interceptor(f RecoveryHandlerFunc) gorpool.Interceptor {
	return func(ctx context.Context, runner gorpool.Job) error {
		defer func() {
			if p := recover(); p != nil {
				f(p)
			}
		}()
		return runner(ctx)
	}
}
