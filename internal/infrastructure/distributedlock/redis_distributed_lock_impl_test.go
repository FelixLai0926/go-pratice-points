package distributedlock

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestAcquireSuccess(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	locker := NewRedisLocker(client)
	lockDuration := 5 * time.Second
	retryInterval := 100 * time.Millisecond

	_, err = locker.Acquire(context.Background(), "test-key", lockDuration, retryInterval)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	if !s.Exists("test-key") {
		t.Errorf("expected key 'test-key' to exist in redis")
	}

	ttl := s.TTL("test-key")
	if err != nil {
		t.Fatalf("failed to get TTL: %v", err)
	}
	if ttl > lockDuration || ttl < lockDuration-1*time.Second {
		t.Errorf("unexpected TTL: got %v, expected around %v", ttl, lockDuration)
	}
}

func TestAcquireFailure(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	s.Set("test-key", "other-token")

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	locker := NewRedisLocker(client)
	lockDuration := 5 * time.Second
	retryInterval := 100 * time.Millisecond

	_, err = locker.Acquire(context.Background(), "test-key", lockDuration, retryInterval)
	if err == nil {
		t.Fatalf("expected Acquire to fail due to key already locked, but it succeeded")
	}
}

func TestReleaseSuccess(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	locker := NewRedisLocker(client)
	lockDuration := 5 * time.Second
	retryInterval := 100 * time.Millisecond

	lock, err := locker.Acquire(context.Background(), "test-key", lockDuration, retryInterval)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	if err := lock.Release(context.Background()); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	if s.Exists("test-key") {
		t.Errorf("expected key 'test-key' to be removed after Release")
	}
}

func TestReleaseFailure(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	locker := NewRedisLocker(client)
	lockDuration := 5 * time.Second
	retryInterval := 100 * time.Millisecond

	lock, err := locker.Acquire(context.Background(), "test-key", lockDuration, retryInterval)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	s.Set("test-key", "different-token")

	err = lock.Release(context.Background())
	if err == nil {
		t.Fatalf("expected Release to fail due to token mismatch")
	}
}

func TestRenewSuccess(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	locker := NewRedisLocker(client)
	lockDuration := 5 * time.Second
	retryInterval := 100 * time.Millisecond

	lock, err := locker.Acquire(context.Background(), "test-key", lockDuration, retryInterval)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	newTTL := 10 * time.Second
	if err := lock.Renew(context.Background(), newTTL); err != nil {
		t.Fatalf("Renew failed: %v", err)
	}

	ttl := s.TTL("test-key")
	if err != nil {
		t.Fatalf("failed to get TTL: %v", err)
	}
	if ttl > newTTL || ttl < newTTL-1*time.Second {
		t.Errorf("unexpected TTL after Renew: got %v, expected around %v", ttl, newTTL)
	}
}

func TestRenewFailure(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	locker := NewRedisLocker(client)
	lockDuration := 5 * time.Second
	retryInterval := 100 * time.Millisecond

	lock, err := locker.Acquire(context.Background(), "test-key", lockDuration, retryInterval)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	s.Set("test-key", "different-token")

	newTTL := 10 * time.Second
	err = lock.Renew(context.Background(), newTTL)
	if err == nil {
		t.Fatalf("expected Renew to fail due to token mismatch")
	}
}
