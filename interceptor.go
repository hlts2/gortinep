package gortinep

import "context"

// Job defines the handler of job for goroutine pool.
type Job func(context.Context) error

// Interceptor provides a hook to intercept the execution of Job.
type Interceptor func(context.Context, Job) error
