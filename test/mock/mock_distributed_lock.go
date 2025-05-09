// Code generated by MockGen. DO NOT EDIT.
// Source: D:/Practice/go-practice/points/internal/domain/distributed_lock.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	domain "points/internal/domain"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockLocker is a mock of Locker interface.
type MockLocker struct {
	ctrl     *gomock.Controller
	recorder *MockLockerMockRecorder
}

// MockLockerMockRecorder is the mock recorder for MockLocker.
type MockLockerMockRecorder struct {
	mock *MockLocker
}

// NewMockLocker creates a new mock instance.
func NewMockLocker(ctrl *gomock.Controller) *MockLocker {
	mock := &MockLocker{ctrl: ctrl}
	mock.recorder = &MockLockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocker) EXPECT() *MockLockerMockRecorder {
	return m.recorder
}

// Acquire mocks base method.
func (m *MockLocker) Acquire(ctx context.Context, key string, lockDuration, retryInterval time.Duration) (domain.Lock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Acquire", ctx, key, lockDuration, retryInterval)
	ret0, _ := ret[0].(domain.Lock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Acquire indicates an expected call of Acquire.
func (mr *MockLockerMockRecorder) Acquire(ctx, key, lockDuration, retryInterval interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acquire", reflect.TypeOf((*MockLocker)(nil).Acquire), ctx, key, lockDuration, retryInterval)
}

// MockLock is a mock of Lock interface.
type MockLock struct {
	ctrl     *gomock.Controller
	recorder *MockLockMockRecorder
}

// MockLockMockRecorder is the mock recorder for MockLock.
type MockLockMockRecorder struct {
	mock *MockLock
}

// NewMockLock creates a new mock instance.
func NewMockLock(ctrl *gomock.Controller) *MockLock {
	mock := &MockLock{ctrl: ctrl}
	mock.recorder = &MockLockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLock) EXPECT() *MockLockMockRecorder {
	return m.recorder
}

// Release mocks base method.
func (m *MockLock) Release(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Release", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Release indicates an expected call of Release.
func (mr *MockLockMockRecorder) Release(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockLock)(nil).Release), ctx)
}

// Renew mocks base method.
func (m *MockLock) Renew(ctx context.Context, ttl time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Renew", ctx, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// Renew indicates an expected call of Renew.
func (mr *MockLockMockRecorder) Renew(ctx, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Renew", reflect.TypeOf((*MockLock)(nil).Renew), ctx, ttl)
}
