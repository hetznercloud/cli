// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: DatacenterClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// MockDatacenterClient is a mock of DatacenterClient interface.
type MockDatacenterClient struct {
	ctrl     *gomock.Controller
	recorder *MockDatacenterClientMockRecorder
}

// MockDatacenterClientMockRecorder is the mock recorder for MockDatacenterClient.
type MockDatacenterClientMockRecorder struct {
	mock *MockDatacenterClient
}

// NewMockDatacenterClient creates a new mock instance.
func NewMockDatacenterClient(ctrl *gomock.Controller) *MockDatacenterClient {
	mock := &MockDatacenterClient{ctrl: ctrl}
	mock.recorder = &MockDatacenterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatacenterClient) EXPECT() *MockDatacenterClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockDatacenterClient) All(arg0 context.Context) ([]*hcloud.Datacenter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Datacenter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockDatacenterClientMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockDatacenterClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockDatacenterClient) AllWithOpts(arg0 context.Context, arg1 hcloud.DatacenterListOpts) ([]*hcloud.Datacenter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Datacenter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockDatacenterClientMockRecorder) AllWithOpts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockDatacenterClient)(nil).AllWithOpts), arg0, arg1)
}

// Get mocks base method.
func (m *MockDatacenterClient) Get(arg0 context.Context, arg1 string) (*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockDatacenterClientMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDatacenterClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockDatacenterClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDatacenterClientMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDatacenterClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *MockDatacenterClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockDatacenterClientMockRecorder) GetByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockDatacenterClient)(nil).GetByName), arg0, arg1)
}

// List mocks base method.
func (m *MockDatacenterClient) List(arg0 context.Context, arg1 hcloud.DatacenterListOpts) ([]*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockDatacenterClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDatacenterClient)(nil).List), arg0, arg1)
}

// Names mocks base method.
func (m *MockDatacenterClient) Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Names indicates an expected call of Names.
func (mr *MockDatacenterClientMockRecorder) Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Names", reflect.TypeOf((*MockDatacenterClient)(nil).Names))
}
