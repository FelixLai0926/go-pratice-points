// Code generated by MockGen. DO NOT EDIT.
// Source: D:/Practice/go-practice/points/internal/domain/port/settings_manager.go

// Package mock is a generated GoMock package.
package mock

import (
	port "points/internal/domain/port"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSettingsManager is a mock of SettingsManager interface.
type MockSettingsManager struct {
	ctrl     *gomock.Controller
	recorder *MockSettingsManagerMockRecorder
}

// MockSettingsManagerMockRecorder is the mock recorder for MockSettingsManager.
type MockSettingsManagerMockRecorder struct {
	mock *MockSettingsManager
}

// NewMockSettingsManager creates a new mock instance.
func NewMockSettingsManager(ctrl *gomock.Controller) *MockSettingsManager {
	mock := &MockSettingsManager{ctrl: ctrl}
	mock.recorder = &MockSettingsManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSettingsManager) EXPECT() *MockSettingsManagerMockRecorder {
	return m.recorder
}

// GetInt mocks base method.
func (m *MockSettingsManager) GetInt(key string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt", key)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetInt indicates an expected call of GetInt.
func (mr *MockSettingsManagerMockRecorder) GetInt(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt", reflect.TypeOf((*MockSettingsManager)(nil).GetInt), key)
}

// GetString mocks base method.
func (m *MockSettingsManager) GetString(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetString indicates an expected call of GetString.
func (mr *MockSettingsManagerMockRecorder) GetString(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockSettingsManager)(nil).GetString), key)
}

// SetDefault mocks base method.
func (m *MockSettingsManager) SetDefault(key string, value interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDefault", key, value)
}

// SetDefault indicates an expected call of SetDefault.
func (mr *MockSettingsManagerMockRecorder) SetDefault(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefault", reflect.TypeOf((*MockSettingsManager)(nil).SetDefault), key, value)
}

// Sub mocks base method.
func (m *MockSettingsManager) Sub(key string) port.SettingsManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sub", key)
	ret0, _ := ret[0].(port.SettingsManager)
	return ret0
}

// Sub indicates an expected call of Sub.
func (mr *MockSettingsManagerMockRecorder) Sub(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sub", reflect.TypeOf((*MockSettingsManager)(nil).Sub), key)
}

// Unmarshal mocks base method.
func (m *MockSettingsManager) Unmarshal(out interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockSettingsManagerMockRecorder) Unmarshal(out interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockSettingsManager)(nil).Unmarshal), out)
}
