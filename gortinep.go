package gortinep

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

const (
	// DefaultPoolSize is default pool size.
	DefaultPoolSize = 128

	// DefaultJobSize is default job size.
	DefaultJobSize = 256
)

// Gortinep is base gortinep interface.
type Gortinep interface {
	Add(Job)
	Start(context.Context) Gortinep
	Stop() Gortinep
	Wait() chan error
}

type (
	gortinep struct {
		running     bool
		poolSize    int
		jobSize     int
		workers     []*worker
		workerWg    *sync.WaitGroup
		jobWg       *sync.WaitGroup
		jobCh       chan Job
		sigDoneCh   chan struct{}
		jobError    *jobError
		interceptor Interceptor
		cancel      context.CancelFunc
	}

	worker struct {
		gp      *gortinep
		running int64
	}

	jobError struct {
		ch     chan error
		closed bool
		mu     *sync.Mutex
	}
)

func newDefaultJobError() *jobError {
	return &jobError{
		ch:     make(chan error),
		closed: false,
		mu:     new(sync.Mutex),
	}
}

func (e *jobError) close() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return
	}

	close(e.ch)
	e.closed = true
}

func (e *jobError) open() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.closed {
		return
	}

	e.ch = make(chan error, cap(e.ch))
	e.closed = false
}

// New creates Gortinep(*gortinep) instance.
func New(opts ...Option) Gortinep {
	gp := newDefaultGortinep()

	for _, opt := range opts {
		opt(gp)
	}

	for i := 0; i < gp.poolSize; i++ {
		gp.workers[i] = newDefaultWorker(gp)
	}

	return gp
}

func newDefaultGortinep() *gortinep {
	return &gortinep{
		running:   false,
		poolSize:  DefaultPoolSize,
		jobSize:   DefaultJobSize,
		workers:   make([]*worker, DefaultPoolSize),
		workerWg:  new(sync.WaitGroup),
		jobWg:     new(sync.WaitGroup),
		jobCh:     make(chan Job, DefaultJobSize),
		sigDoneCh: make(chan struct{}),
		jobError:  newDefaultJobError(),
	}
}

func newDefaultWorker(gp *gortinep) *worker {
	return &worker{
		gp:      gp,
		running: 0,
	}
}

// Start starts all goroutine pool with context.
func (gp *gortinep) Start(ctx context.Context) Gortinep {
	if gp.running {
		return gp
	}

	ctx, gp.cancel = context.WithCancel(ctx)

	go gp.watchShutdownSignal(gp.sigDoneCh, gp.cancel)

	for _, worker := range gp.workers {
		if atomic.LoadInt64(&worker.running) == 0 {
			gp.workerWg.Add(1)
			go worker.start(ctx)
		}
	}

	gp.running = true

	return gp
}

func (gp *gortinep) watchShutdownSignal(sigDoneCh chan struct{}, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(sigCh)
		close(sigCh)
	}()

	for {
		select {
		case <-gp.sigDoneCh:
			return
		case <-sigCh:
			cancel()
			gp.workerWg.Wait()
		}
	}
}

// Stop stops all goroutine pool.
// If job is being executed in goroutine pool, wait until it is finished and stop the groutine pool.
func (gp *gortinep) Stop() Gortinep {
	if !gp.running {
		return gp
	}

	gp.cancel()

	gp.running = false

	gp.sigDoneCh <- struct{}{}

	return gp
}

// Add adds job into gorutine pool. job is processed asynchronously.
func (gp *gortinep) Add(job Job) {
	if gp.jobError != nil {
		gp.jobError.open()
	}

	gp.jobWg.Add(1)
	gp.jobCh <- job
}

// Wait return error channel for job error processed by goroutine worker.
// If the error channel is not set, wait for all jobs to end and return.
func (gp *gortinep) Wait() chan error {
	if gp.jobError == nil {
		gp.jobWg.Wait()
		return nil
	}

	go func() {
		gp.jobWg.Wait()
		gp.jobError.close()
	}()

	return gp.jobError.ch
}

func (w *worker) start(ctx context.Context) {
	atomic.AddInt64(&w.running, 1)

	defer func() {
		atomic.StoreInt64(&w.running, 0)
		w.gp.workerWg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case j := <-w.gp.jobCh:

			// Send job error to channel.
			// If error channel is nil, do nothing.
			w.sendJobError(w.execute(ctx, j))
			w.gp.jobWg.Done()
		}
	}
}

func (w *worker) execute(ctx context.Context, job Job) (err error) {
	if w.gp.interceptor == nil {
		err = job(ctx)
	} else {
		err = w.gp.interceptor(ctx, job)
	}

	return
}

func (w *worker) sendJobError(err error) {
	if w.gp.jobError == nil {
		return
	}

	w.gp.jobError.ch <- err
}
