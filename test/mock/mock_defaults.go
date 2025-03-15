// Code generated by MockGen. DO NOT EDIT.
// Source: D:/Practice/go-practice/points/internal/domain/port/defaults.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDefaults is a mock of Defaults interface.
type MockDefaults struct {
	ctrl     *gomock.Controller
	recorder *MockDefaultsMockRecorder
}

// MockDefaultsMockRecorder is the mock recorder for MockDefaults.
type MockDefaultsMockRecorder struct {
	mock *MockDefaults
}

// NewMockDefaults creates a new mock instance.
func NewMockDefaults(ctrl *gomock.Controller) *MockDefaults {
	mock := &MockDefaults{ctrl: ctrl}
	mock.recorder = &MockDefaultsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDefaults) EXPECT() *MockDefaultsMockRecorder {
	return m.recorder
}

// Set mocks base method.
func (m *MockDefaults) Set(ptr interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ptr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockDefaultsMockRecorder) Set(ptr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockDefaults)(nil).Set), ptr)
}
