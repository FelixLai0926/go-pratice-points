package domain

import (
	"context"
	"time"
)

type Locker interface {
	Acquire(ctx context.Context, key string, lockDuration, retryInterval time.Duration) (Lock, error)
}

type Lock interface {
	Release(ctx context.Context) error
	Renew(ctx context.Context, ttl time.Duration) error
}
