package grpool

import (
	"context"
	"fmt"
	"testing"
)

func testInterceptor(ctx context.Context, runner Runner) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover")
		}
	}()

	return runner(ctx)
}

func TestNew(t *testing.T) {
	g := New(
		WithPoolSize(100),
		WithUnaryInterceptor(testInterceptor),
	)

	g.Sync(context.Background(), func(ctx context.Context) error {
		return nil
	})
}
