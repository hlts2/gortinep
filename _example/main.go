package main

import (
	"context"
	"fmt"

	"github.com/hlts2/gortinep"
	"github.com/hlts2/gortinep/middlewares"
	"github.com/hlts2/gortinep/middlewares/logger/zap"
	"github.com/hlts2/gortinep/middlewares/recovery"
	"go.uber.org/zap"
)

func main() {
	z := zap.NewExample()

	g := gortinep.New(
		gortinep.WithErrorChannel(make(chan error, 1)),
		gortinep.WithWorkerSize(256),
		gortinep.WithInterceptor(
			middlewares.ChainInterceptors(
				gortinep_recovery.Interceptor(
					func(p interface{}) {
						z.Info("recovery from panic")
					},
				),
				gortinep_zap.Interceptor(
					z,
				),
			),
		),
	).Start(context.Background())
	defer g.Stop()

	const jobSize = 100000

	for i := 0; i < jobSize; i++ {
		// Register job.
		g.Add(func(context.Context) error {
			z.Info("finish job")
			return nil
		})
	}

	// Get error of job. If execution of all jobs is completed, exit Wait function.
	for err := range g.Wait() {
		if err != nil {
			fmt.Println(err)
		}
	}

	// Register job again.
	for i := 0; i < jobSize; i++ {
		g.Add(func(context.Context) error {
			z.Info("finish job")
			return nil
		})
	}

	for err := range g.Wait() {
		if err != nil {
			fmt.Println(err)
		}
	}
}
