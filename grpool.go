package grpool

import (
	"context"
	"sync"
)

// DefaultPoolSize --
const DefaultPoolSize = 100

// GrPool --
type GrPool interface {
	GetCurrentPoolSize() int
	Start(context.Context) GrPool
	Stop() GrPool
}

type grPool struct {
	running     bool
	poolSize    int
	workers     []*worker
	runnerCh    chan Runner
	interceptor Interceptor
}

type worker struct {
	gp      *grPool
	mu      *sync.Mutex
	running bool
	killCh  chan struct{}
}

// New --
func New(opts ...Option) GrPool {
	gp := createDefaultGrpool()

	for _, opt := range opts {
		opt(gp)
	}

	for i := 0; i < gp.poolSize; i++ {
		gp.workers = append(gp.workers, createDefaultWorker(gp))
	}

	return gp
}

func createDefaultGrpool() *grPool {
	return &grPool{
		running:  false,
		poolSize: DefaultPoolSize,
		workers:  make([]*worker, 0, DefaultPoolSize),
		runnerCh: make(chan Runner),
	}
}

func createDefaultWorker(gp *grPool) *worker {
	return &worker{
		gp:      gp,
		mu:      new(sync.Mutex),
		running: false,
		killCh:  make(chan struct{}),
	}
}

func (gp *grPool) Start(ctx context.Context) GrPool {
	for _, worker := range gp.workers {
		if !worker.running {
			// start worker
			worker.running = true
			worker.start(ctx)
		}
	}
	gp.running = true
	return gp
}

func (gp *grPool) Stop() GrPool {
	for _, worker := range gp.workers {
		if worker.running {
			worker.killCh <- struct{}{}
			worker.running = false
		}
	}

	gp.running = false
	return gp
}

func (gp *grPool) GetCurrentPoolSize() int {
	return len(gp.workers)
}

func (gp *grPool) Add(ctx context.Context, runner Runner) {
	gp.runnerCh <- runner
}

func (w *worker) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-w.killCh:
				return
			case <-ctx.Done():
				w.mu.Lock()
				w.running = false
				w.mu.Unlock()
				return
			}
		}
	}()
}

func (w *worker) execute(ctx context.Context, runner Runner) {
	if w.gp.interceptor == nil {
		runner(ctx)
	} else {
		w.gp.interceptor(ctx, runner)
	}
}
