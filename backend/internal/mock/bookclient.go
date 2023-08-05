// Code generated by MockGen. DO NOT EDIT.
// Source: ./interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBookClient is a mock of BookClient interface.
type MockBookClient struct {
	ctrl     *gomock.Controller
	recorder *MockBookClientMockRecorder
}

// MockBookClientMockRecorder is the mock recorder for MockBookClient.
type MockBookClientMockRecorder struct {
	mock *MockBookClient
}

// NewMockBookClient creates a new mock instance.
func NewMockBookClient(ctrl *gomock.Controller) *MockBookClient {
	mock := &MockBookClient{ctrl: ctrl}
	mock.recorder = &MockBookClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBookClient) EXPECT() *MockBookClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockBookClient) Get(ctx context.Context, url string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, url)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockBookClientMockRecorder) Get(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBookClient)(nil).Get), ctx, url)
}
