package distributedlock

import (
	"context"
	"points/pkg/module/test"
	"testing"
	"time"
)

func TestRedisLockClient_AcquireAndRelease(t *testing.T) {
	mr, client := test.NewDummyRedis(t)
	defer mr.Close()

	lockClient := NewRedisLockClient(client)
	ctx := context.Background()
	key := "test-lock"
	ttl := 5 * time.Second

	lock, err := lockClient.Acquire(ctx, key, ttl)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	if !mr.Exists(key) {
		t.Errorf("lock key %q not found in Redis after acquire", key)
	}

	if err := lockClient.Release(ctx, lock); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	if mr.Exists(key) {
		t.Errorf("lock key %q still exists in Redis after release", key)
	}
}

func TestRedisLockClient_Renew(t *testing.T) {
	mr, client := test.NewDummyRedis(t)
	defer mr.Close()

	lockClient := NewRedisLockClient(client)
	ctx := context.Background()
	key := "renew-lock"
	ttl := 2 * time.Second

	lock, err := lockClient.Acquire(ctx, key, ttl)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	time.Sleep(1 * time.Second)

	if err := lockClient.Renew(ctx, lock, 4*time.Second); err != nil {
		t.Fatalf("Renew failed: %v", err)
	}

	remainingTTL := mr.TTL(key)
	if remainingTTL < 3*time.Second {
		t.Errorf("expected renewed TTL to be at least 3s, got %v", remainingTTL)
	}

	if err := lockClient.Release(ctx, lock); err != nil {
		t.Fatalf("Release failed: %v", err)
	}
}
