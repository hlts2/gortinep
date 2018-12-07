package grpool

import "context"

type Runner func(ctx context.Context) error
