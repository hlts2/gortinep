package grpool

import (
	"context"
	"sync"
)

// DefaultPoolSize is default pool size
const DefaultPoolSize = 100

// GrPool is base grpool interface
type GrPool interface {
	Add(Runner)
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

// New create GrPool(*grPool) instance
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

// Start starts all goroutine pool with context
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

// Stop stops all goroutine pool
// If job is being executed in goroutine pool, wait until it is finished and stop the groutine pool
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

// GetCurrentPoolSize returns current goroutine pool size
func (gp *grPool) GetCurrentPoolSize() int {
	size := 0

	for _, worker := range gp.workers {
		worker.mu.Lock()
		size++
		worker.mu.Unlock()
	}
	return size
}

// Add adds job into gorutine pool. job is processed asynchronously
func (gp *grPool) Add(runner Runner) {
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
			case r := <-w.gp.runnerCh:
				w.execute(r)
			}
		}
	}()
}

func (w *worker) execute(runner Runner) {
	if w.gp.interceptor == nil {
		runner()
	} else {
		w.gp.interceptor(runner)
	}
}
