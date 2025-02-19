package distributedlock

import (
	"context"
	"testing"
	"time"

	"points/pkg/module/test"

	"github.com/stretchr/testify/assert"
)

func TestRedisLockClient_AcquireAndRelease(t *testing.T) {
	mr, client := test.NewDummyRedis(t)
	defer mr.Close()

	lockClient := NewRedisLockClient(client)
	ctx := context.Background()
	key := "test-lock"
	ttl := 5 * time.Second

	lock, err := lockClient.Acquire(ctx, key, ttl)
	assert.NoError(t, err, "Acquire should succeed")
	assert.NotNil(t, lock, "Lock should not be nil")

	assert.True(t, mr.Exists(key), "lock key %q should exist in Redis after acquire", key)

	err = lockClient.Release(ctx, lock)
	assert.NoError(t, err, "Release should succeed")

	assert.False(t, mr.Exists(key), "lock key %q should not exist in Redis after release", key)
}

func TestRedisLockClient_Renew(t *testing.T) {
	mr, client := test.NewDummyRedis(t)
	defer mr.Close()

	lockClient := NewRedisLockClient(client)
	ctx := context.Background()
	key := "renew-lock"
	ttl := 2 * time.Second

	lock, err := lockClient.Acquire(ctx, key, ttl)
	assert.NoError(t, err, "Acquire should succeed")
	assert.NotNil(t, lock, "Lock should not be nil")

	time.Sleep(1 * time.Second)

	err = lockClient.Renew(ctx, lock, 4*time.Second)
	assert.NoError(t, err, "Renew should succeed")

	remainingTTL := mr.TTL(key)
	assert.GreaterOrEqual(t, remainingTTL, 3*time.Second, "Expected renewed TTL to be at least 3s, got %v", remainingTTL)

	err = lockClient.Release(ctx, lock)
	assert.NoError(t, err, "Release should succeed")
}
