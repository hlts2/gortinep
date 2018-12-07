package grpool

import "context"

// Handler --
type Handler func(context.Context)

// Interceptor --
type Interceptor func(context.Context, Handler)
