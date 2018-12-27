package gortinep

import (
	"context"
	"os"
	"os/signal"
	"sync"
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

type gortinep struct {
	running      bool
	poolSize     int
	jobSize      int
	workers      []*worker
	wjobg        *sync.WaitGroup
	jobCh        chan Job
	wrapperrCh   *wrapperrCh
	sigDoneCh    chan struct{}
	workerDoneCh chan struct{}
	closedErrCh  bool
	mu           *sync.Mutex

	interceptor Interceptor
}

type worker struct {
	gp      *gortinep
	killCh  chan struct{}
	running bool
}

type wrapperrCh struct {
	ch     chan error
	closed bool
}

func (wech *wrapperrCh) reopen() {
	if !wech.closed {
		return
	}

	wech.ch = make(chan error, cap(wech.ch))
	wech.closed = false
}

func (wech *wrapperrCh) close() {
	if wech.closed {
		return
	}

	close(wech.ch)
	wech.closed = true
}

// New creates Gortinep(*gortinep) instance.
func New(opts ...Option) Gortinep {
	gp := createDefaultGortinep()

	for _, opt := range opts {
		opt(gp)
	}

	for i := 0; i < gp.poolSize; i++ {
		gp.workers[i] = createDefaultWorker(gp)
	}

	return gp
}

func createDefaultGortinep() *gortinep {
	return &gortinep{
		running:      false,
		poolSize:     DefaultPoolSize,
		jobSize:      DefaultJobSize,
		workers:      make([]*worker, DefaultPoolSize),
		wjobg:        new(sync.WaitGroup),
		jobCh:        make(chan Job, DefaultJobSize),
		sigDoneCh:    make(chan struct{}),
		workerDoneCh: make(chan struct{}),
		mu:           new(sync.Mutex),
	}
}

func createDefaultWorker(gp *gortinep) *worker {
	return &worker{
		gp:      gp,
		running: false,
		killCh:  make(chan struct{}),
	}
}

// Start starts all goroutine pool with context.
func (gp *gortinep) Start(ctx context.Context) Gortinep {
	if gp.running {
		return gp
	}

	// starts os signal observer.
	cctx := gp.signalObserver(ctx, gp.sigDoneCh)

	for _, worker := range gp.workers {
		if !worker.running {
			// starts worker with context.
			worker.start(cctx)
			worker.running = true
		}
	}

	gp.running = true
	return gp
}

// Stop stops all goroutine pool.
// If job is being executed in goroutine pool, wait until it is finished and stop the groutine pool.
func (gp *gortinep) Stop() Gortinep {
	if !gp.running {
		return gp
	}

	for _, worker := range gp.workers {
		if worker.running {
			worker.killCh <- struct{}{}
			worker.running = false
		}
	}

	// stops os signal observer.
	gp.sigDoneCh <- struct{}{}

	gp.running = false
	return gp
}

func (gp *gortinep) signalObserver(ctx context.Context, sigDoneCh chan struct{}) context.Context {
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
				gp.waitWorkers()
			case <-sigDoneCh:
				return
			}
		}
	}()

	return cctx
}

// waitWorkers waits for all workers to finish.
func (gp *gortinep) waitWorkers() {
	defer func() {
		close(gp.workerDoneCh)
		close(gp.sigDoneCh)
	}()

	n := 0
	for _ = range gp.workerDoneCh {
		n++
		if n == len(gp.workers) {
			return
		}
	}
}

// Add adds job into gorutine pool. job is processed asynchronously.
func (gp *gortinep) Add(job Job) {
	if gp.wrapperrCh != nil {
		gp.mu.Lock()
		gp.wrapperrCh.reopen()
		gp.mu.Unlock()
	}
	gp.wjobg.Add(1)
	gp.jobCh <- job
}

// Wait return error channel for job error processed by goroutine worker.
// If the error channel is not set, wait for all jobs to end and return.
func (gp *gortinep) Wait() chan error {
	if gp.wrapperrCh == nil {
		gp.wjobg.Wait()
		return nil
	}

	go func() {
		gp.wjobg.Wait()
		gp.mu.Lock()
		gp.wrapperrCh.close()
		gp.mu.Unlock()
	}()
	return gp.wrapperrCh.ch
}

func (w *worker) start(ctx context.Context) {
	go func() {
		defer func() {
			w.gp.workerDoneCh <- struct{}{}
		}()
		for {
			select {
			case <-w.killCh:
				return
			case <-ctx.Done():
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
	if w.gp.wrapperrCh == nil {
		return
	}

	w.gp.wrapperrCh.ch <- err
}
