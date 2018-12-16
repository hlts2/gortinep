package grpool

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func signalObserver(ctx context.Context) context.Context {
	sigCh := make(chan os.Signal, 1)
	cctx, cancel := context.WithCancel(ctx)

	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)

	go func() {
		for {
			select {
			case <-sigCh:
				cancel()
			}
		}
	}()

	return cctx
}
