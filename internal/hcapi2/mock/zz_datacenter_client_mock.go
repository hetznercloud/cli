// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: DatacenterClient)
//
// Generated by this command:
//
//	mockgen -package mock -destination zz_datacenter_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 DatacenterClient
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// MockDatacenterClient is a mock of DatacenterClient interface.
type MockDatacenterClient struct {
	ctrl     *gomock.Controller
	recorder *MockDatacenterClientMockRecorder
	isgomock struct{}
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
func (m *MockDatacenterClient) All(ctx context.Context) ([]*hcloud.Datacenter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", ctx)
	ret0, _ := ret[0].([]*hcloud.Datacenter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockDatacenterClientMockRecorder) All(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockDatacenterClient)(nil).All), ctx)
}

// AllWithOpts mocks base method.
func (m *MockDatacenterClient) AllWithOpts(ctx context.Context, opts hcloud.DatacenterListOpts) ([]*hcloud.Datacenter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", ctx, opts)
	ret0, _ := ret[0].([]*hcloud.Datacenter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockDatacenterClientMockRecorder) AllWithOpts(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockDatacenterClient)(nil).AllWithOpts), ctx, opts)
}

// Get mocks base method.
func (m *MockDatacenterClient) Get(ctx context.Context, idOrName string) (*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, idOrName)
	ret0, _ := ret[0].(*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockDatacenterClientMockRecorder) Get(ctx, idOrName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDatacenterClient)(nil).Get), ctx, idOrName)
}

// GetByID mocks base method.
func (m *MockDatacenterClient) GetByID(ctx context.Context, id int64) (*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDatacenterClientMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDatacenterClient)(nil).GetByID), ctx, id)
}

// GetByName mocks base method.
func (m *MockDatacenterClient) GetByName(ctx context.Context, name string) (*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name)
	ret0, _ := ret[0].(*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockDatacenterClientMockRecorder) GetByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockDatacenterClient)(nil).GetByName), ctx, name)
}

// List mocks base method.
func (m *MockDatacenterClient) List(ctx context.Context, opts hcloud.DatacenterListOpts) ([]*hcloud.Datacenter, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, opts)
	ret0, _ := ret[0].([]*hcloud.Datacenter)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockDatacenterClientMockRecorder) List(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDatacenterClient)(nil).List), ctx, opts)
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
