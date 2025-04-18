// Code generated by MockGen. DO NOT EDIT.
// Source: D:/Practice/go-practice/points/internal/domain/port/copier.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCopier is a mock of Copier interface.
type MockCopier struct {
	ctrl     *gomock.Controller
	recorder *MockCopierMockRecorder
}

// MockCopierMockRecorder is the mock recorder for MockCopier.
type MockCopierMockRecorder struct {
	mock *MockCopier
}

// NewMockCopier creates a new mock instance.
func NewMockCopier(ctrl *gomock.Controller) *MockCopier {
	mock := &MockCopier{ctrl: ctrl}
	mock.recorder = &MockCopierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCopier) EXPECT() *MockCopierMockRecorder {
	return m.recorder
}

// Copy mocks base method.
func (m *MockCopier) Copy(dst, src interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy", dst, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// Copy indicates an expected call of Copy.
func (mr *MockCopierMockRecorder) Copy(dst, src interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockCopier)(nil).Copy), dst, src)
}
