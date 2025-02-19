package distributedlock

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"go.uber.org/zap"
)

type LockClient interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error)
	Release(ctx context.Context, lock *redislock.Lock) error
	Renew(ctx context.Context, lock *redislock.Lock, ttl time.Duration) error
}

func WithLock(ctx context.Context, client LockClient, key string, ttl time.Duration, operation func() error) error {
	lock, err := client.Acquire(ctx, key, ttl)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer func() {
		if err := client.Release(ctx, lock); err != nil {
			zap.L().Error("failed to release lock", zap.Error(err))
		}
	}()
	return operation()
}
