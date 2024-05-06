// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: ActionClient)

// Package hcapi2_mock is a generated GoMock package.
package hcapi2_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// MockActionClient is a mock of ActionClient interface.
type MockActionClient struct {
	ctrl     *gomock.Controller
	recorder *MockActionClientMockRecorder
}

// MockActionClientMockRecorder is the mock recorder for MockActionClient.
type MockActionClientMockRecorder struct {
	mock *MockActionClient
}

// NewMockActionClient creates a new mock instance.
func NewMockActionClient(ctrl *gomock.Controller) *MockActionClient {
	mock := &MockActionClient{ctrl: ctrl}
	mock.recorder = &MockActionClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActionClient) EXPECT() *MockActionClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockActionClient) All(arg0 context.Context) ([]*hcloud.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockActionClientMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockActionClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockActionClient) AllWithOpts(arg0 context.Context, arg1 hcloud.ActionListOpts) ([]*hcloud.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockActionClientMockRecorder) AllWithOpts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockActionClient)(nil).AllWithOpts), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockActionClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockActionClientMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockActionClient)(nil).GetByID), arg0, arg1)
}

// List mocks base method.
func (m *MockActionClient) List(arg0 context.Context, arg1 hcloud.ActionListOpts) ([]*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockActionClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockActionClient)(nil).List), arg0, arg1)
}

// WaitFor mocks base method.
func (m *MockActionClient) WaitFor(arg0 context.Context, arg1 ...*hcloud.Action) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WaitFor", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitFor indicates an expected call of WaitFor.
func (mr *MockActionClientMockRecorder) WaitFor(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitFor", reflect.TypeOf((*MockActionClient)(nil).WaitFor), varargs...)
}

// WaitForFunc mocks base method.
func (m *MockActionClient) WaitForFunc(arg0 context.Context, arg1 func(*hcloud.Action) error, arg2 ...*hcloud.Action) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WaitForFunc", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitForFunc indicates an expected call of WaitForFunc.
func (mr *MockActionClientMockRecorder) WaitForFunc(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForFunc", reflect.TypeOf((*MockActionClient)(nil).WaitForFunc), varargs...)
}

// WatchOverallProgress mocks base method.
func (m *MockActionClient) WatchOverallProgress(arg0 context.Context, arg1 []*hcloud.Action) (<-chan int, <-chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchOverallProgress", arg0, arg1)
	ret0, _ := ret[0].(<-chan int)
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// WatchOverallProgress indicates an expected call of WatchOverallProgress.
func (mr *MockActionClientMockRecorder) WatchOverallProgress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchOverallProgress", reflect.TypeOf((*MockActionClient)(nil).WatchOverallProgress), arg0, arg1)
}

// WatchProgress mocks base method.
func (m *MockActionClient) WatchProgress(arg0 context.Context, arg1 *hcloud.Action) (<-chan int, <-chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchProgress", arg0, arg1)
	ret0, _ := ret[0].(<-chan int)
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// WatchProgress indicates an expected call of WatchProgress.
func (mr *MockActionClientMockRecorder) WatchProgress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchProgress", reflect.TypeOf((*MockActionClient)(nil).WatchProgress), arg0, arg1)
}
