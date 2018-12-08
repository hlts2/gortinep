package grpool

import (
	"context"
	"fmt"
)

func testInterceptor(ctx context.Context, runner Runner) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover")
		}
	}()

	return runner(ctx)
}
