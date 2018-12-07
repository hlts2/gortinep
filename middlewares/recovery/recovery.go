package recovery

import (
	"context"

	"github.com/hlts2/grpool"
)

// Recovery --
type Recovery func(ctx context.Context)

// UnaryInterceptor --
func UnaryInterceptor(rcv Recovery) grpool.Interceptor {
	return func(ctx context.Context, runner grpool.Runner) {
		defer rcv(ctx)
		runner(ctx)
	}
}
