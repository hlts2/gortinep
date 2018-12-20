package middlewares

import (
	"context"
	"testing"

	"github.com/hlts2/grpool"
)

func TestChainInterceptors(t *testing.T) {
	t.Run("n == 1", func(t *testing.T) {
		var got string
		chainInterceptor := ChainInterceptors(
			func(ctx context.Context, job grpool.Job) error {
				got += "enter interceptor1\n"
				err := job(ctx)
				got += "finish interceptor1\n"
				return err
			},
			func(ctx context.Context, job grpool.Job) error {
				got += "enter interceptor2\n"
				err := job(ctx)
				got += "finish interceptor2\n"
				return err
			},
		)

		err := chainInterceptor(context.Background(), func(ctx context.Context) error {
			got += "job\n"
			return nil
		})

		if err != nil {
			t.Errorf("chainInterceptor is err: %v", err)
		}

		expected := "enter interceptor1\nenter interceptor2\njob\nfinish interceptor2\nfinish interceptor1\n"
		if got != expected {
			t.Errorf("chainInterceptor is wrong. expected %v, but got: %v", expected, got)
		}

	})
}
