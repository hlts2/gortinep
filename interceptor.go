package grpool

import "context"

// Interceptor --
type Interceptor func(context.Context, Runner)
