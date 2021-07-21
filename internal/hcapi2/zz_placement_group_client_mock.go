// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: PlacementGroupClient)

// Package hcapi2 is a generated GoMock package.
package hcapi2

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hcloud "github.com/hetznercloud/hcloud-go/hcloud"
)

// MockPlacementGroupClient is a mock of PlacementGroupClient interface.
type MockPlacementGroupClient struct {
	ctrl     *gomock.Controller
	recorder *MockPlacementGroupClientMockRecorder
}

// MockPlacementGroupClientMockRecorder is the mock recorder for MockPlacementGroupClient.
type MockPlacementGroupClientMockRecorder struct {
	mock *MockPlacementGroupClient
}

// NewMockPlacementGroupClient creates a new mock instance.
func NewMockPlacementGroupClient(ctrl *gomock.Controller) *MockPlacementGroupClient {
	mock := &MockPlacementGroupClient{ctrl: ctrl}
	mock.recorder = &MockPlacementGroupClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlacementGroupClient) EXPECT() *MockPlacementGroupClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockPlacementGroupClient) All(arg0 context.Context) ([]*hcloud.PlacementGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.PlacementGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockPlacementGroupClientMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockPlacementGroupClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockPlacementGroupClient) AllWithOpts(arg0 context.Context, arg1 hcloud.PlacementGroupListOpts) ([]*hcloud.PlacementGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.PlacementGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockPlacementGroupClientMockRecorder) AllWithOpts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockPlacementGroupClient)(nil).AllWithOpts), arg0, arg1)
}

// Delete mocks base method.
func (m *MockPlacementGroupClient) Delete(arg0 context.Context, arg1 *hcloud.PlacementGroup) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockPlacementGroupClientMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPlacementGroupClient)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockPlacementGroupClient) Get(arg0 context.Context, arg1 string) (*hcloud.PlacementGroup, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PlacementGroup)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockPlacementGroupClientMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPlacementGroupClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockPlacementGroupClient) GetByID(arg0 context.Context, arg1 int) (*hcloud.PlacementGroup, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PlacementGroup)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockPlacementGroupClientMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockPlacementGroupClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *MockPlacementGroupClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.PlacementGroup, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PlacementGroup)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockPlacementGroupClientMockRecorder) GetByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockPlacementGroupClient)(nil).GetByName), arg0, arg1)
}

// LabelKeys mocks base method.
func (m *MockPlacementGroupClient) LabelKeys(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LabelKeys", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// LabelKeys indicates an expected call of LabelKeys.
func (mr *MockPlacementGroupClientMockRecorder) LabelKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LabelKeys", reflect.TypeOf((*MockPlacementGroupClient)(nil).LabelKeys), arg0)
}

// List mocks base method.
func (m *MockPlacementGroupClient) List(arg0 context.Context, arg1 hcloud.PlacementGroupListOpts) ([]*hcloud.PlacementGroup, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.PlacementGroup)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockPlacementGroupClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPlacementGroupClient)(nil).List), arg0, arg1)
}

// Names mocks base method.
func (m *MockPlacementGroupClient) Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Names indicates an expected call of Names.
func (mr *MockPlacementGroupClientMockRecorder) Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Names", reflect.TypeOf((*MockPlacementGroupClient)(nil).Names))
}

// Update mocks base method.
func (m *MockPlacementGroupClient) Update(arg0 context.Context, arg1 *hcloud.PlacementGroup, arg2 hcloud.PlacementGroupUpdateOpts) (*hcloud.PlacementGroup, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.PlacementGroup)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockPlacementGroupClientMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPlacementGroupClient)(nil).Update), arg0, arg1, arg2)
}
