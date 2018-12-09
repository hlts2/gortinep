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

func (gr *grPool) Start(ctx context.Context) GrPool {
	for _, worker := range gr.workers {
		if !worker.running {
			// start worker
			worker.running = true
			worker.start(ctx)
		}
	}
	gr.running = true
	return gr
}

func (gr *grPool) Stop() GrPool {
	for _, worker := range gr.workers {
		if worker.running {
			worker.killCh <- struct{}{}
			worker.running = false
		}
	}

	gr.running = false
	return gr
}

func (gr *grPool) GetCurrentPoolSize() int {
	return len(gr.workers)
}

func (gr *grPool) Add(runner Runner) {
	gr.runnerCh <- runner
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
			case r := <-w.gp.runnerCh:
				w.execute(ctx, r)
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
