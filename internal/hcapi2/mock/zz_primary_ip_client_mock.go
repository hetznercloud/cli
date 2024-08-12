// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: PrimaryIPClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// MockPrimaryIPClient is a mock of PrimaryIPClient interface.
type MockPrimaryIPClient struct {
	ctrl     *gomock.Controller
	recorder *MockPrimaryIPClientMockRecorder
}

// MockPrimaryIPClientMockRecorder is the mock recorder for MockPrimaryIPClient.
type MockPrimaryIPClientMockRecorder struct {
	mock *MockPrimaryIPClient
}

// NewMockPrimaryIPClient creates a new mock instance.
func NewMockPrimaryIPClient(ctrl *gomock.Controller) *MockPrimaryIPClient {
	mock := &MockPrimaryIPClient{ctrl: ctrl}
	mock.recorder = &MockPrimaryIPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrimaryIPClient) EXPECT() *MockPrimaryIPClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockPrimaryIPClient) All(arg0 context.Context) ([]*hcloud.PrimaryIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.PrimaryIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockPrimaryIPClientMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockPrimaryIPClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockPrimaryIPClient) AllWithOpts(arg0 context.Context, arg1 hcloud.PrimaryIPListOpts) ([]*hcloud.PrimaryIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.PrimaryIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockPrimaryIPClientMockRecorder) AllWithOpts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockPrimaryIPClient)(nil).AllWithOpts), arg0, arg1)
}

// Assign mocks base method.
func (m *MockPrimaryIPClient) Assign(arg0 context.Context, arg1 hcloud.PrimaryIPAssignOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Assign", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Assign indicates an expected call of Assign.
func (mr *MockPrimaryIPClientMockRecorder) Assign(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Assign", reflect.TypeOf((*MockPrimaryIPClient)(nil).Assign), arg0, arg1)
}

// ChangeDNSPtr mocks base method.
func (m *MockPrimaryIPClient) ChangeDNSPtr(arg0 context.Context, arg1 hcloud.PrimaryIPChangeDNSPtrOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDNSPtr", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeDNSPtr indicates an expected call of ChangeDNSPtr.
func (mr *MockPrimaryIPClientMockRecorder) ChangeDNSPtr(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDNSPtr", reflect.TypeOf((*MockPrimaryIPClient)(nil).ChangeDNSPtr), arg0, arg1)
}

// ChangeProtection mocks base method.
func (m *MockPrimaryIPClient) ChangeProtection(arg0 context.Context, arg1 hcloud.PrimaryIPChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeProtection", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeProtection indicates an expected call of ChangeProtection.
func (mr *MockPrimaryIPClientMockRecorder) ChangeProtection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeProtection", reflect.TypeOf((*MockPrimaryIPClient)(nil).ChangeProtection), arg0, arg1)
}

// Create mocks base method.
func (m *MockPrimaryIPClient) Create(arg0 context.Context, arg1 hcloud.PrimaryIPCreateOpts) (*hcloud.PrimaryIPCreateResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PrimaryIPCreateResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockPrimaryIPClientMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPrimaryIPClient)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockPrimaryIPClient) Delete(arg0 context.Context, arg1 *hcloud.PrimaryIP) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockPrimaryIPClientMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPrimaryIPClient)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockPrimaryIPClient) Get(arg0 context.Context, arg1 string) (*hcloud.PrimaryIP, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PrimaryIP)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockPrimaryIPClientMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPrimaryIPClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockPrimaryIPClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.PrimaryIP, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PrimaryIP)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockPrimaryIPClientMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockPrimaryIPClient)(nil).GetByID), arg0, arg1)
}

// GetByIP mocks base method.
func (m *MockPrimaryIPClient) GetByIP(arg0 context.Context, arg1 string) (*hcloud.PrimaryIP, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIP", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PrimaryIP)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByIP indicates an expected call of GetByIP.
func (mr *MockPrimaryIPClientMockRecorder) GetByIP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIP", reflect.TypeOf((*MockPrimaryIPClient)(nil).GetByIP), arg0, arg1)
}

// GetByName mocks base method.
func (m *MockPrimaryIPClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.PrimaryIP, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.PrimaryIP)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockPrimaryIPClientMockRecorder) GetByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockPrimaryIPClient)(nil).GetByName), arg0, arg1)
}

// IPv4Names mocks base method.
func (m *MockPrimaryIPClient) IPv4Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IPv4Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// IPv4Names indicates an expected call of IPv4Names.
func (mr *MockPrimaryIPClientMockRecorder) IPv4Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IPv4Names", reflect.TypeOf((*MockPrimaryIPClient)(nil).IPv4Names))
}

// IPv6Names mocks base method.
func (m *MockPrimaryIPClient) IPv6Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IPv6Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// IPv6Names indicates an expected call of IPv6Names.
func (mr *MockPrimaryIPClientMockRecorder) IPv6Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IPv6Names", reflect.TypeOf((*MockPrimaryIPClient)(nil).IPv6Names))
}

// LabelKeys mocks base method.
func (m *MockPrimaryIPClient) LabelKeys(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LabelKeys", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// LabelKeys indicates an expected call of LabelKeys.
func (mr *MockPrimaryIPClientMockRecorder) LabelKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LabelKeys", reflect.TypeOf((*MockPrimaryIPClient)(nil).LabelKeys), arg0)
}

// List mocks base method.
func (m *MockPrimaryIPClient) List(arg0 context.Context, arg1 hcloud.PrimaryIPListOpts) ([]*hcloud.PrimaryIP, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.PrimaryIP)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockPrimaryIPClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPrimaryIPClient)(nil).List), arg0, arg1)
}

// Names mocks base method.
func (m *MockPrimaryIPClient) Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Names indicates an expected call of Names.
func (mr *MockPrimaryIPClientMockRecorder) Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Names", reflect.TypeOf((*MockPrimaryIPClient)(nil).Names))
}

// Unassign mocks base method.
func (m *MockPrimaryIPClient) Unassign(arg0 context.Context, arg1 int64) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unassign", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Unassign indicates an expected call of Unassign.
func (mr *MockPrimaryIPClientMockRecorder) Unassign(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unassign", reflect.TypeOf((*MockPrimaryIPClient)(nil).Unassign), arg0, arg1)
}

// Update mocks base method.
func (m *MockPrimaryIPClient) Update(arg0 context.Context, arg1 *hcloud.PrimaryIP, arg2 hcloud.PrimaryIPUpdateOpts) (*hcloud.PrimaryIP, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.PrimaryIP)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockPrimaryIPClientMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPrimaryIPClient)(nil).Update), arg0, arg1, arg2)
}
