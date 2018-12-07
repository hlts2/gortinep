package grpool

import (
	"testing"
)

func TestNew(t *testing.T) {
	_ = NewGrPool(
		WithGrPoolSize(100),
		WithUnaryInterceptor(
			ChainInterceptors(
				UnalyRecoveryInterceptor(),
			),
		),
	)
}
