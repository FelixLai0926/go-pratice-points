package locking

import (
	"context"
	"errors"
	"testing"
	"time"

	"points/internal/shared/apperror"
	"points/internal/shared/errcode"
	"points/test/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWithTradeLock_Success_GoMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLock := mock.NewMockLock(ctrl)
	mockLock.EXPECT().Release(gomock.Any()).Return(nil).AnyTimes()
	mockLock.EXPECT().Renew(gomock.Any(), gomock.Any()).Return(nil).MinTimes(1)

	mockLocker := mock.NewMockLocker(ctrl)
	mockLocker.EXPECT().Acquire(gomock.Any(), "test-key", gomock.Any(), gomock.Any()).Return(mockLock, nil)

	svc := &accountLockApplicationService{
		locker:        mockLocker,
		lockDuration:  100 * time.Millisecond,
		retryInterval: 50 * time.Millisecond,
	}

	opCalled := false
	operation := func() error {
		opCalled = true
		time.Sleep(150 * time.Millisecond)
		return nil
	}

	err := svc.WithTradeLock(context.Background(), "test-key", operation)
	assert.NoError(t, err)
	assert.True(t, opCalled, "operation should be called")
}

func TestWithTradeLock_AcquireFailure_GoMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLocker := mock.NewMockLocker(ctrl)
	acquireErr := errors.New("acquire error")
	mockLocker.EXPECT().Acquire(gomock.Any(), "test-key", gomock.Any(), gomock.Any()).Return(nil, acquireErr)

	svc := &accountLockApplicationService{
		locker:        mockLocker,
		lockDuration:  100 * time.Millisecond,
		retryInterval: 50 * time.Millisecond,
	}

	operation := func() error { return nil }
	err := svc.WithTradeLock(context.Background(), "test-key", operation)
	assert.Error(t, err)
	appErr, ok := err.(*apperror.AppError)
	assert.True(t, ok, "err should be AppError")
	assert.Equal(t, appErr.Code, errcode.ErrDistrubutedLockAcquire)
}

func TestWithTradeLock_RenewFailure_GoMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLock := mock.NewMockLock(ctrl)
	renewErr := errors.New("renew error")
	mockLock.EXPECT().Release(gomock.Any()).Return(nil).AnyTimes()
	mockLock.EXPECT().Renew(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockLock.EXPECT().Renew(gomock.Any(), gomock.Any()).Return(renewErr)

	mockLocker := mock.NewMockLocker(ctrl)
	mockLocker.EXPECT().Acquire(gomock.Any(), "test-key", gomock.Any(), gomock.Any()).Return(mockLock, nil)

	svc := &accountLockApplicationService{
		locker:        mockLocker,
		lockDuration:  100 * time.Millisecond,
		retryInterval: 50 * time.Millisecond,
	}

	operation := func() error {
		time.Sleep(300 * time.Millisecond)
		return nil
	}

	err := svc.WithTradeLock(context.Background(), "test-key", operation)
	assert.Error(t, err)
	appErr, ok := err.(*apperror.AppError)
	assert.True(t, ok, "err should be AppError")
	assert.Equal(t, appErr.Code, errcode.ErrDistrubutedLockRenew)
}
