package recovery

import (
	"github.com/hlts2/grpool"
)

// Recovery --
type Recovery func()

// UnaryInterceptor --
func UnaryInterceptor(rcv Recovery) grpool.Interceptor {
	return func(runner grpool.Runner) error {
		defer rcv()
		return runner()
	}
}
