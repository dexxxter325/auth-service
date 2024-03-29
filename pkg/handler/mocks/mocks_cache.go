// Code generated by MockGen. DO NOT EDIT.
// Source: C:/Users/prsok/GolandProjects/CRUD_API/pkg/handler/cache.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockCaches is a mock of Caches interface.
type MockCaches struct {
	ctrl     *gomock.Controller
	recorder *MockCachesMockRecorder
}

// MockCachesMockRecorder is the mock recorder for MockCaches.
type MockCachesMockRecorder struct {
	mock *MockCaches
}

// NewMockCaches creates a new mock instance.
func NewMockCaches(ctrl *gomock.Controller) *MockCaches {
	mock := &MockCaches{ctrl: ctrl}
	mock.recorder = &MockCachesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCaches) EXPECT() *MockCachesMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockCaches) Get(key string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockCachesMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCaches)(nil).Get), key)
}

// Set mocks base method.
func (m *MockCaches) Set(cacheKey string, data interface{}, duration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", cacheKey, data, duration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCachesMockRecorder) Set(cacheKey, data, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCaches)(nil).Set), cacheKey, data, duration)
}
