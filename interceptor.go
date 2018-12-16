package grpool

// Runner defines the handler of goruotine pool task
type Runner func() error

// Interceptor provides a hook to intercept the execution of Runner
type Interceptor func(Runner) error
