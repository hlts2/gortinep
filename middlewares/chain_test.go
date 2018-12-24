package middlewares

import (
	"context"
	"testing"

	"github.com/hlts2/gortinep"
)

func TestChainInterceptors(t *testing.T) {
	t.Run("The number of interceptors is greater than 1", func(t *testing.T) {
		var got string
		chainInterceptor := ChainInterceptors(
			func(ctx context.Context, job gortinep.Job) error {
				got += "enter first interceptor\n"
				err := job(ctx)
				got += "finish first interceptor\n"
				return err
			},
			func(ctx context.Context, job gortinep.Job) error {
				got += "enter second interceptor\n"
				err := job(ctx)
				got += "finish second interceptor\n"
				return err
			},
		)

		err := chainInterceptor(context.Background(), func(ctx context.Context) error {
			got += "job\n"
			return nil
		})

		if err != nil {
			t.Errorf("ChainInterceptor is err: %v", err)
		}

		expected := "enter first interceptor\nenter second interceptor\njob\nfinish second interceptor\nfinish first interceptor\n"
		if got != expected {
			t.Errorf("ChainInterceptor is wrong. expected %v, but got: %v", expected, got)
		}
	})

	t.Run("The number of interceptors is 1", func(t *testing.T) {
		var got string
		chainInterceptor := ChainInterceptors(
			func(ctx context.Context, job gortinep.Job) error {
				got += "enter first interceptor\n"
				err := job(ctx)
				got += "finish first interceptor\n"
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

		expected := "enter first interceptor\njob\nfinish first interceptor\n"
		if got != expected {
			t.Errorf("ChainInterceptor is wrong. expected %v, but got: %v", expected, got)
		}
	})

	t.Run("The number of interceptors is 0", func(t *testing.T) {
		var got string
		chainInterceptor := ChainInterceptors()

		err := chainInterceptor(context.Background(), func(ctx context.Context) error {
			got += "job\n"
			return nil
		})

		if err != nil {
			t.Errorf("ChainInterceptor is err: %v", err)
		}

		expected := "job\n"
		if got != expected {
			t.Errorf("ChainInterceptor is wrong. expected %v, but got: %v", expected, got)
		}
	})
}
