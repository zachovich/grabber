package grabber

import "context"

type RateLimiter interface {
	WainN(ctx context.Context, n int) (err error)
}
