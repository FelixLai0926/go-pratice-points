package distributedlock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type redisLockClient struct {
	client *redis.Client
	locker *redislock.Client
}

func NewRedisLockClient(client *redis.Client) LockClient {
	return &redisLockClient{
		client: client,
		locker: redislock.New(client),
	}
}

func (r *redisLockClient) Acquire(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	opts := &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(100 * time.Millisecond),
	}

	return r.locker.Obtain(ctx, key, ttl, opts)
}

func (r *redisLockClient) Release(ctx context.Context, lock *redislock.Lock) error {
	return lock.Release(ctx)
}

func (r *redisLockClient) Renew(ctx context.Context, lock *redislock.Lock, ttl time.Duration) error {
	return lock.Refresh(ctx, ttl, nil)
}
