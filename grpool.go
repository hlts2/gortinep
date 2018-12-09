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

type runnable struct {
	ctx    context.Context
	runner Runner
}

type grPool struct {
	running     bool
	poolSize    int
	workers     []*worker
	runnableCh  chan runnable
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
		running:    false,
		poolSize:   DefaultPoolSize,
		workers:    make([]*worker, 0, DefaultPoolSize),
		runnableCh: make(chan runnable),
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

func (gr *grPool) Add(ctx context.Context, runner Runner) {
	gr.runnableCh <- runnable{
		ctx:    ctx,
		runner: runner,
	}
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
			case r := <-w.gp.runnableCh:
				w.execute(r.ctx, r.runner)
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
