// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: SSHKeyClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// MockSSHKeyClient is a mock of SSHKeyClient interface.
type MockSSHKeyClient struct {
	ctrl     *gomock.Controller
	recorder *MockSSHKeyClientMockRecorder
}

// MockSSHKeyClientMockRecorder is the mock recorder for MockSSHKeyClient.
type MockSSHKeyClientMockRecorder struct {
	mock *MockSSHKeyClient
}

// NewMockSSHKeyClient creates a new mock instance.
func NewMockSSHKeyClient(ctrl *gomock.Controller) *MockSSHKeyClient {
	mock := &MockSSHKeyClient{ctrl: ctrl}
	mock.recorder = &MockSSHKeyClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSSHKeyClient) EXPECT() *MockSSHKeyClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockSSHKeyClient) All(arg0 context.Context) ([]*hcloud.SSHKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.SSHKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockSSHKeyClientMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockSSHKeyClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockSSHKeyClient) AllWithOpts(arg0 context.Context, arg1 hcloud.SSHKeyListOpts) ([]*hcloud.SSHKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.SSHKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockSSHKeyClientMockRecorder) AllWithOpts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockSSHKeyClient)(nil).AllWithOpts), arg0, arg1)
}

// Create mocks base method.
func (m *MockSSHKeyClient) Create(arg0 context.Context, arg1 hcloud.SSHKeyCreateOpts) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockSSHKeyClientMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSSHKeyClient)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockSSHKeyClient) Delete(arg0 context.Context, arg1 *hcloud.SSHKey) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockSSHKeyClientMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSSHKeyClient)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockSSHKeyClient) Get(arg0 context.Context, arg1 string) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockSSHKeyClientMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSSHKeyClient)(nil).Get), arg0, arg1)
}

// GetByFingerprint mocks base method.
func (m *MockSSHKeyClient) GetByFingerprint(arg0 context.Context, arg1 string) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFingerprint", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByFingerprint indicates an expected call of GetByFingerprint.
func (mr *MockSSHKeyClientMockRecorder) GetByFingerprint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFingerprint", reflect.TypeOf((*MockSSHKeyClient)(nil).GetByFingerprint), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockSSHKeyClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockSSHKeyClientMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockSSHKeyClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *MockSSHKeyClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockSSHKeyClientMockRecorder) GetByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockSSHKeyClient)(nil).GetByName), arg0, arg1)
}

// LabelKeys mocks base method.
func (m *MockSSHKeyClient) LabelKeys(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LabelKeys", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// LabelKeys indicates an expected call of LabelKeys.
func (mr *MockSSHKeyClientMockRecorder) LabelKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LabelKeys", reflect.TypeOf((*MockSSHKeyClient)(nil).LabelKeys), arg0)
}

// List mocks base method.
func (m *MockSSHKeyClient) List(arg0 context.Context, arg1 hcloud.SSHKeyListOpts) ([]*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockSSHKeyClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSSHKeyClient)(nil).List), arg0, arg1)
}

// Names mocks base method.
func (m *MockSSHKeyClient) Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Names indicates an expected call of Names.
func (mr *MockSSHKeyClientMockRecorder) Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Names", reflect.TypeOf((*MockSSHKeyClient)(nil).Names))
}

// Update mocks base method.
func (m *MockSSHKeyClient) Update(arg0 context.Context, arg1 *hcloud.SSHKey, arg2 hcloud.SSHKeyUpdateOpts) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockSSHKeyClientMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSSHKeyClient)(nil).Update), arg0, arg1, arg2)
}
