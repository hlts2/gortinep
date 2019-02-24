package gortinep

import "sync"

// Option configures gortinep.
type Option func(*gortinep)

// WithWorkerSize returns an option that sets the worker size.
func WithWorkerSize(size int) Option {
	return func(gp *gortinep) {
		if size < 1 {
			return
		}
		gp.workerSize = size
		gp.workers = make([]*worker, gp.workerSize)
	}
}

// WithJobSize returns an option that sets the job size.
func WithJobSize(size int) Option {
	return func(gp *gortinep) {
		if size < 1 {
			return
		}
		gp.jobSize = size
		gp.jobCh = make(chan Job, size)
	}
}

// WithErrorChannel returns an option that sets channel for job error processed by goroutine worker.
// The result of each goroutine is sent to this channe.
func WithErrorChannel(ch chan error) Option {
	return func(gp *gortinep) {
		if ch == nil {
			return
		}

		gp.jobError = &jobError{
			ch:     ch,
			closed: false,
			mu:     new(sync.Mutex),
		}
	}
}

// WithInterceptor returns an option that sets the Interceptor implementation.
func WithInterceptor(interceptor Interceptor) Option {
	return func(gp *gortinep) {
		gp.interceptor = interceptor
	}
}
