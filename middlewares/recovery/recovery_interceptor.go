package grpool_recovery

import (
	"github.com/hlts2/grpool"
)

// RecoveryHandlerFunc is a function that recovers from the panic `p`.
type RecoveryHandlerFunc func(p interface{})

// Interceptor returns a new interceptor for panic recovery.
func Interceptor(f RecoveryHandlerFunc) grpool.Interceptor {
	return func(runner grpool.Job) error {
		defer func() {
			if p := recover(); p != nil {
				f(p)
			}
		}()
		return runner()
	}
}