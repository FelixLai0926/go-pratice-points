package distributedlock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
)

type FailingLockClient struct{}

func NewFailingLockClient() LockClient {
	return &FailingLockClient{}
}

func (c *FailingLockClient) Acquire(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	return nil, errors.New("failed to acquire lock")
}

func (c *FailingLockClient) Release(ctx context.Context, lock *redislock.Lock) error {
	return nil
}

func (c *FailingLockClient) Renew(ctx context.Context, lock *redislock.Lock, ttl time.Duration) error {
	return nil
}
