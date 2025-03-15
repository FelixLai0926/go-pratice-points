package locking

import (
	"context"
	"fmt"
	"points/internal/domain"
	"points/internal/domain/port"
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"
	"time"

	"go.uber.org/zap"
)

type AccountLockApplicationService interface {
	WithAccountTradeLock(ctx context.Context, from, to int64, fn func() error) error
}

type accountLockApplicationService struct {
	locker        domain.Locker
	lockDuration  time.Duration
	retryInterval time.Duration
	config        port.Config
}

func NewAccountLockService(locker domain.Locker, config port.Config) AccountLockApplicationService {
	lockDuration, retryInterval := initTTL(config)
	return &accountLockApplicationService{
		locker:        locker,
		lockDuration:  lockDuration,
		retryInterval: retryInterval,
		config:        config,
	}
}

func (a *accountLockApplicationService) WithAccountTradeLock(ctx context.Context, from, to int64, fn func() error) error {
	lockKey := getLockKey(from, to)
	return a.WithTradeLock(ctx, lockKey, fn)
}

func (a *accountLockApplicationService) WithTradeLock(ctx context.Context, key string, operation func() error) error {
	lock, err := a.locker.Acquire(ctx, key, a.lockDuration, a.retryInterval)
	if err != nil {
		return apperror.Wrap(errcode.ErrDistrubutedLockAcquire, "failed to acquire lock", err)
	}
	defer func() {
		if err := lock.Release(ctx); err != nil {
			zap.L().Error("failed to release lock", zap.Error(err))
		}
	}()

	renewCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		ticker := time.NewTicker(a.lockDuration / 2)
		defer ticker.Stop()

		for {
			select {
			case <-renewCtx.Done():
				return
			case <-ticker.C:
				if err := lock.Renew(renewCtx, a.lockDuration); err != nil {
					errCh <- apperror.Wrap(errcode.ErrDistrubutedLockRenew, "failed to renew lock", err)
					return
				}
			}
		}
	}()

	opCh := make(chan error, 1)
	go func() {
		opCh <- operation()
	}()

	select {
	case opErr := <-opCh:
		select {
		case renewErr := <-errCh:
			return renewErr
		default:
			return opErr
		}
	case renewErr := <-errCh:
		cancel()
		return renewErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

func initTTL(config port.Config) (time.Duration, time.Duration) {
	config.SetDefaultInt("LOCK_DURATION", 5)
	config.SetDefaultInt("RETRY_INTERVAL", 100)

	lockDuration := config.GetInt("LOCK_DURATION")
	retryInterval := config.GetInt("RETRY_INTERVAL")

	return time.Duration(lockDuration) * time.Second, time.Duration(retryInterval) * time.Second
}

func getLockKey(from, to int64) string {
	if from < to {
		return fmt.Sprintf("transfer_lock:%d:%d", from, to)
	}
	return fmt.Sprintf("transfer_lock:%d:%d", to, from)
}
