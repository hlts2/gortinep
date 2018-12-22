package grpool

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// DefaultPoolSize is default pool size.
const DefaultPoolSize = 256

// GrPool is base grpool interface.
type GrPool interface {
	Add(Job)
	Start(context.Context) GrPool
	Stop() GrPool
	Wait() chan error
}

type grPool struct {
	running               bool
	poolSize              int
	workers               []*worker
	wjobg                 *sync.WaitGroup
	jobCh                 chan Job
	errCh                 chan error
	sigDoneCh             chan struct{}
	workerDoneCh          chan struct{}
	isClosedErrCh         bool
	runningWorkerObserver bool
	mu                    *sync.Mutex

	interceptor Interceptor
}

type worker struct {
	gp      *grPool
	running bool
	stopCh  chan struct{}
}

// New creates GrPool(*grPool) instance.
func New(opts ...Option) GrPool {
	gp := createDefaultGrpool()

	for _, opt := range opts {
		opt(gp)
	}

	for i := 0; i < gp.poolSize; i++ {
		gp.workers[i] = createDefaultWorker(gp)
	}

	return gp
}

func createDefaultGrpool() *grPool {
	return &grPool{
		running:      false,
		poolSize:     DefaultPoolSize,
		workers:      make([]*worker, DefaultPoolSize),
		wjobg:        new(sync.WaitGroup),
		jobCh:        make(chan Job),
		sigDoneCh:    make(chan struct{}),
		workerDoneCh: make(chan struct{}),
		mu:           new(sync.Mutex),
	}
}

func createDefaultWorker(gp *grPool) *worker {
	return &worker{
		gp:      gp,
		running: false,
		stopCh:  make(chan struct{}),
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

	for _, worker := range gp.workers {
		if worker.running {
			worker.stopCh <- struct{}{}
		}
	}

	// stops os signal observer.
	gp.sigDoneCh <- struct{}{}

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
		defer func() {
			signal.Stop(sigCh)
			close(sigCh)
		}()

		for {
			select {
			case <-sigCh:
				cancel()
				gp.workerShutdownObserver()
			case <-doneCh:
				return
			}
		}
	}()

	return cctx
}

func (gp *grPool) workerShutdownObserver() {
	defer gp.mu.Unlock()
	gp.mu.Lock()

	if gp.runningWorkerObserver {
		return
	}

	gp.runningWorkerObserver = true

	go func() {
		defer func() {
			close(gp.sigDoneCh)
			close(gp.workerDoneCh)
		}()

		n := 0
		for _ = range gp.workerDoneCh {
			n++
			if n == len(gp.workers) {
				return
			}
		}
	}()
}

// Add adds job into gorutine pool. job is processed asynchronously.
func (gp *grPool) Add(job Job) {
	if gp.errCh != nil {
		gp.mu.Lock()
		if gp.isClosedErrCh {
			gp.isClosedErrCh = false
			gp.errCh = make(chan error, cap(gp.errCh))
		}
		gp.mu.Unlock()
	}
	gp.wjobg.Add(1)
	gp.jobCh <- job
}

// Wait return error channel for job error processed by goroutine worker.
// If the error channel is not set, wait for all jobs to end and return.
func (gp *grPool) Wait() chan error {
	if gp.errCh == nil {
		gp.wjobg.Wait()
		return nil
	}

	go func() {
		gp.wjobg.Wait()
		close(gp.errCh)
		gp.mu.Lock()
		gp.isClosedErrCh = true
		gp.mu.Unlock()
	}()
	return gp.errCh
}

func (w *worker) start(ctx context.Context) {
	w.running = true

	go func() {
		defer func() {
			w.gp.workerDoneCh <- struct{}{}
		}()
		for {
			select {
			case <-w.stopCh:
				w.stop()
				return
			case <-ctx.Done():
				w.stop()
				return
			case j := <-w.gp.jobCh:

				// Notifies job error into error channel.
				// if error channel is nil, do nothing.
				w.notifyJobError(w.execute(ctx, j))
				w.gp.wjobg.Done()
			}
		}
	}()
}

func (w *worker) stop() {
	w.gp.mu.Lock()
	w.running = false
	w.gp.mu.Unlock()
}

func (w *worker) execute(ctx context.Context, job Job) error {
	var err error

	if w.gp.interceptor == nil {
		err = job(ctx)
	} else {
		err = w.gp.interceptor(ctx, job)
	}

	return err
}

func (w *worker) notifyJobError(err error) {
	if w.gp.errCh == nil {
		return
	}

	w.gp.errCh <- err
}
