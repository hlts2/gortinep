package grpool

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// DefaultPoolSize is default pool size.
const DefaultPoolSize = 100

// GrPool is base grpool interface.
type GrPool interface {
	Add(Job)
	GetCurrentPoolSize() int
	Start(context.Context) GrPool
	Stop() GrPool
	Error() chan error
}

type grPool struct {
	running     bool
	poolSize    int
	workers     []*worker
	jobCh       chan Job
	errCh       chan error
	sigDoneCh   chan struct{}
	interceptor Interceptor
}

type worker struct {
	gp      *grPool
	mu      *sync.Mutex
	running bool
	killCh  chan struct{}
}

// New creates GrPool(*grPool) instance.
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
		running:   false,
		poolSize:  DefaultPoolSize,
		workers:   make([]*worker, 0, DefaultPoolSize),
		jobCh:     make(chan Job),
		sigDoneCh: make(chan struct{}),
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

// Start starts all goroutine pool with context.
func (gp *grPool) Start(ctx context.Context) GrPool {
	if gp.running {
		return gp
	}

	// starts os signal observer.
	cctx := gp.signalObserver(ctx, gp.sigDoneCh)

	for _, worker := range gp.workers {
		if !worker.running {
			// starts worker.
			worker.running = true
			worker.start(cctx)
		}
	}
	gp.running = true
	return gp
}

// Stop stops all goroutine pool.
// If job is being executed in goroutine pool, wait until it is finished and stop the groutine pool.
func (gp *grPool) Stop() GrPool {
	if !gp.running {
		return gp
	}

	// stops os signal observer.
	gp.sigDoneCh <- struct{}{}

	for _, worker := range gp.workers {
		if worker.running {
			worker.killCh <- struct{}{}
			worker.running = false
		}
	}

	gp.running = false
	return gp
}

func (gp *grPool) signalObserver(ctx context.Context, doneCh chan struct{}) context.Context {
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
			case <-doneCh:
				return
			}
		}
	}()

	return cctx
}

// GetCurrentPoolSize returns current goroutine pool size.
func (gp *grPool) GetCurrentPoolSize() int {
	size := 0

	for _, worker := range gp.workers {
		worker.mu.Lock()
		size++
		worker.mu.Unlock()
	}
	return size
}

// Add adds job into gorutine pool. job is processed asynchronously.
func (gp *grPool) Add(job Job) {
	gp.jobCh <- job
}

func (gp *grPool) Error() chan error {
	if gp.errCh == nil {
		return nil
	}
	return gp.errCh
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
			case j := <-w.gp.jobCh:

				// Notifies job error into error channel.
				// if error channel is nil, do nothing.
				w.notifyJobError(w.execute(j))
			}
		}
	}()
}

func (w *worker) execute(job Job) error {
	var err error

	if w.gp.interceptor == nil {
		err = job()
	} else {
		err = w.gp.interceptor(job)
	}

	return err
}

func (w *worker) notifyJobError(err error) {
	if w.gp.errCh == nil {
		return
	}

	w.gp.errCh <- err
}
