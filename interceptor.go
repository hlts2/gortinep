package grpool

// Job defines the handler of job for goruotine pool
type Job func() error

// Interceptor provides a hook to intercept the execution of Job
type Interceptor func(Job) error
