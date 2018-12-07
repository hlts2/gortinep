package grpool

import (
	"context"
	"fmt"
)

// UnalyRecoveryInterceptor --
func UnalyRecoveryInterceptor() Interceptor {
	return func(ctx context.Context, runner Runner) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovery")
			}
		}()

		runner(ctx)
	}
}
