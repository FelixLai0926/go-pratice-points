package distributedlock

import (
	"context"
	"fmt"
	"points/internal/domain"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type RedisLocker struct {
	client *redis.Client
}

var _ domain.Locker = (*RedisLocker)(nil)

func NewRedisLocker(client *redis.Client) domain.Locker {
	return &RedisLocker{
		client: client,
	}
}

func (r *RedisLocker) Acquire(ctx context.Context, key string, lockDuration, retryInterval time.Duration) (domain.Lock, error) {
	opts := &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(retryInterval),
	}

	lock, err := redislock.New(r.client).Obtain(ctx, key, lockDuration, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return &RedisLock{
		lock: lock,
	}, nil
}

type RedisLock struct {
	lock *redislock.Lock
}

var _ domain.Lock = (*RedisLock)(nil)

func (l *RedisLock) Release(ctx context.Context) error {
	return l.lock.Release(ctx)
}

func (l *RedisLock) Renew(ctx context.Context, ttl time.Duration) error {
	return l.lock.Refresh(ctx, ttl, nil)
}
