package grpool

import (
	"context"
)

// DefaultPoolSize --
const DefaultPoolSize = 100

// GrPool --
type GrPool interface {
	Start() GrPool
	Stop() GrPool
}

// Runnable --
type runnable struct {
	runner Runner
	ctx    context.Context
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
		running: false,
		killCh:  make(chan struct{}),
	}
}

func (gr *grPool) Start() GrPool {
	for _, worker := range gr.workers {
		if !worker.running {
			// start worker
			worker.running = true
			go worker.start()
		}
	}
	gr.running = true
	return gr
}

func (gr *grPool) Stop() GrPool {
	for _, worker := range gr.workers {
		worker.killCh <- struct{}{}
		worker.running = false
	}

	gr.running = false
	return gr
}

func (gr *grPool) Add(ctx context.Context, runner Runner) {
	gr.runnableCh <- runnable{
		runner: runner,
		ctx:    ctx,
	}
}

func (w *worker) start() {
	for {
		select {
		case <-w.killCh:
			return
		case r := <-w.gp.runnableCh:
			w.execute(r)
		}
	}
}

func (w *worker) execute(r runnable) {
	if w.gp.interceptor == nil {
		r.runner(r.ctx)
	} else {
		w.gp.interceptor(r.ctx, r.runner)
	}
}
