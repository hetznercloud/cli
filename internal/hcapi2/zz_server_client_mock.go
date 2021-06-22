// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: ServerClient)

// Package hcapi2 is a generated GoMock package.
package hcapi2

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hcloud "github.com/hetznercloud/hcloud-go/hcloud"
)

// MockServerClient is a mock of ServerClient interface.
type MockServerClient struct {
	ctrl     *gomock.Controller
	recorder *MockServerClientMockRecorder
}

// MockServerClientMockRecorder is the mock recorder for MockServerClient.
type MockServerClientMockRecorder struct {
	mock *MockServerClient
}

// NewMockServerClient creates a new mock instance.
func NewMockServerClient(ctrl *gomock.Controller) *MockServerClient {
	mock := &MockServerClient{ctrl: ctrl}
	mock.recorder = &MockServerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerClient) EXPECT() *MockServerClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockServerClient) All(arg0 context.Context) ([]*hcloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockServerClientMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockServerClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockServerClient) AllWithOpts(arg0 context.Context, arg1 hcloud.ServerListOpts) ([]*hcloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockServerClientMockRecorder) AllWithOpts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockServerClient)(nil).AllWithOpts), arg0, arg1)
}

// AttachISO mocks base method.
func (m *MockServerClient) AttachISO(arg0 context.Context, arg1 *hcloud.Server, arg2 *hcloud.ISO) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachISO", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachISO indicates an expected call of AttachISO.
func (mr *MockServerClientMockRecorder) AttachISO(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachISO", reflect.TypeOf((*MockServerClient)(nil).AttachISO), arg0, arg1, arg2)
}

// AttachToNetwork mocks base method.
func (m *MockServerClient) AttachToNetwork(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerAttachToNetworkOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachToNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachToNetwork indicates an expected call of AttachToNetwork.
func (mr *MockServerClientMockRecorder) AttachToNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachToNetwork", reflect.TypeOf((*MockServerClient)(nil).AttachToNetwork), arg0, arg1, arg2)
}

// ChangeAliasIPs mocks base method.
func (m *MockServerClient) ChangeAliasIPs(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerChangeAliasIPsOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAliasIPs", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeAliasIPs indicates an expected call of ChangeAliasIPs.
func (mr *MockServerClientMockRecorder) ChangeAliasIPs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAliasIPs", reflect.TypeOf((*MockServerClient)(nil).ChangeAliasIPs), arg0, arg1, arg2)
}

// ChangeDNSPtr mocks base method.
func (m *MockServerClient) ChangeDNSPtr(arg0 context.Context, arg1 *hcloud.Server, arg2 string, arg3 *string) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDNSPtr", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeDNSPtr indicates an expected call of ChangeDNSPtr.
func (mr *MockServerClientMockRecorder) ChangeDNSPtr(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDNSPtr", reflect.TypeOf((*MockServerClient)(nil).ChangeDNSPtr), arg0, arg1, arg2, arg3)
}

// ChangeProtection mocks base method.
func (m *MockServerClient) ChangeProtection(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeProtection", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeProtection indicates an expected call of ChangeProtection.
func (mr *MockServerClientMockRecorder) ChangeProtection(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeProtection", reflect.TypeOf((*MockServerClient)(nil).ChangeProtection), arg0, arg1, arg2)
}

// ChangeType mocks base method.
func (m *MockServerClient) ChangeType(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerChangeTypeOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeType", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeType indicates an expected call of ChangeType.
func (mr *MockServerClientMockRecorder) ChangeType(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeType", reflect.TypeOf((*MockServerClient)(nil).ChangeType), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *MockServerClient) Create(arg0 context.Context, arg1 hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(hcloud.ServerCreateResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockServerClientMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockServerClient)(nil).Create), arg0, arg1)
}

// CreateImage mocks base method.
func (m *MockServerClient) CreateImage(arg0 context.Context, arg1 *hcloud.Server, arg2 *hcloud.ServerCreateImageOpts) (hcloud.ServerCreateImageResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImage", arg0, arg1, arg2)
	ret0, _ := ret[0].(hcloud.ServerCreateImageResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateImage indicates an expected call of CreateImage.
func (mr *MockServerClientMockRecorder) CreateImage(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImage", reflect.TypeOf((*MockServerClient)(nil).CreateImage), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockServerClient) Delete(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockServerClientMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockServerClient)(nil).Delete), arg0, arg1)
}

// DetachFromNetwork mocks base method.
func (m *MockServerClient) DetachFromNetwork(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerDetachFromNetworkOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachFromNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DetachFromNetwork indicates an expected call of DetachFromNetwork.
func (mr *MockServerClientMockRecorder) DetachFromNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachFromNetwork", reflect.TypeOf((*MockServerClient)(nil).DetachFromNetwork), arg0, arg1, arg2)
}

// DetachISO mocks base method.
func (m *MockServerClient) DetachISO(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachISO", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DetachISO indicates an expected call of DetachISO.
func (mr *MockServerClientMockRecorder) DetachISO(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachISO", reflect.TypeOf((*MockServerClient)(nil).DetachISO), arg0, arg1)
}

// DisableBackup mocks base method.
func (m *MockServerClient) DisableBackup(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableBackup", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DisableBackup indicates an expected call of DisableBackup.
func (mr *MockServerClientMockRecorder) DisableBackup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableBackup", reflect.TypeOf((*MockServerClient)(nil).DisableBackup), arg0, arg1)
}

// DisableRescue mocks base method.
func (m *MockServerClient) DisableRescue(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableRescue", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DisableRescue indicates an expected call of DisableRescue.
func (mr *MockServerClientMockRecorder) DisableRescue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableRescue", reflect.TypeOf((*MockServerClient)(nil).DisableRescue), arg0, arg1)
}

// EnableBackup mocks base method.
func (m *MockServerClient) EnableBackup(arg0 context.Context, arg1 *hcloud.Server, arg2 string) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableBackup", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EnableBackup indicates an expected call of EnableBackup.
func (mr *MockServerClientMockRecorder) EnableBackup(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableBackup", reflect.TypeOf((*MockServerClient)(nil).EnableBackup), arg0, arg1, arg2)
}

// EnableRescue mocks base method.
func (m *MockServerClient) EnableRescue(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerEnableRescueOpts) (hcloud.ServerEnableRescueResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableRescue", arg0, arg1, arg2)
	ret0, _ := ret[0].(hcloud.ServerEnableRescueResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EnableRescue indicates an expected call of EnableRescue.
func (mr *MockServerClientMockRecorder) EnableRescue(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableRescue", reflect.TypeOf((*MockServerClient)(nil).EnableRescue), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockServerClient) Get(arg0 context.Context, arg1 string) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockServerClientMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockServerClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockServerClient) GetByID(arg0 context.Context, arg1 int) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockServerClientMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockServerClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *MockServerClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockServerClientMockRecorder) GetByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockServerClient)(nil).GetByName), arg0, arg1)
}

// GetMetrics mocks base method.
func (m *MockServerClient) GetMetrics(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerGetMetricsOpts) (*hcloud.ServerMetrics, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.ServerMetrics)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockServerClientMockRecorder) GetMetrics(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockServerClient)(nil).GetMetrics), arg0, arg1, arg2)
}

// LabelKeys mocks base method.
func (m *MockServerClient) LabelKeys(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LabelKeys", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// LabelKeys indicates an expected call of LabelKeys.
func (mr *MockServerClientMockRecorder) LabelKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LabelKeys", reflect.TypeOf((*MockServerClient)(nil).LabelKeys), arg0)
}

// List mocks base method.
func (m *MockServerClient) List(arg0 context.Context, arg1 hcloud.ServerListOpts) ([]*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockServerClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockServerClient)(nil).List), arg0, arg1)
}

// Names mocks base method.
func (m *MockServerClient) Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Names indicates an expected call of Names.
func (mr *MockServerClientMockRecorder) Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Names", reflect.TypeOf((*MockServerClient)(nil).Names))
}

// Poweroff mocks base method.
func (m *MockServerClient) Poweroff(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poweroff", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Poweroff indicates an expected call of Poweroff.
func (mr *MockServerClientMockRecorder) Poweroff(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poweroff", reflect.TypeOf((*MockServerClient)(nil).Poweroff), arg0, arg1)
}

// Poweron mocks base method.
func (m *MockServerClient) Poweron(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poweron", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Poweron indicates an expected call of Poweron.
func (mr *MockServerClientMockRecorder) Poweron(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poweron", reflect.TypeOf((*MockServerClient)(nil).Poweron), arg0, arg1)
}

// Reboot mocks base method.
func (m *MockServerClient) Reboot(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reboot", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reboot indicates an expected call of Reboot.
func (mr *MockServerClientMockRecorder) Reboot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reboot", reflect.TypeOf((*MockServerClient)(nil).Reboot), arg0, arg1)
}

// Rebuild mocks base method.
func (m *MockServerClient) Rebuild(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerRebuildOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rebuild", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Rebuild indicates an expected call of Rebuild.
func (mr *MockServerClientMockRecorder) Rebuild(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rebuild", reflect.TypeOf((*MockServerClient)(nil).Rebuild), arg0, arg1, arg2)
}

// RequestConsole mocks base method.
func (m *MockServerClient) RequestConsole(arg0 context.Context, arg1 *hcloud.Server) (hcloud.ServerRequestConsoleResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestConsole", arg0, arg1)
	ret0, _ := ret[0].(hcloud.ServerRequestConsoleResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RequestConsole indicates an expected call of RequestConsole.
func (mr *MockServerClientMockRecorder) RequestConsole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestConsole", reflect.TypeOf((*MockServerClient)(nil).RequestConsole), arg0, arg1)
}

// Reset mocks base method.
func (m *MockServerClient) Reset(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reset indicates an expected call of Reset.
func (mr *MockServerClientMockRecorder) Reset(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockServerClient)(nil).Reset), arg0, arg1)
}

// ResetPassword mocks base method.
func (m *MockServerClient) ResetPassword(arg0 context.Context, arg1 *hcloud.Server) (hcloud.ServerResetPasswordResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", arg0, arg1)
	ret0, _ := ret[0].(hcloud.ServerResetPasswordResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockServerClientMockRecorder) ResetPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockServerClient)(nil).ResetPassword), arg0, arg1)
}

// ServerName mocks base method.
func (m *MockServerClient) ServerName(arg0 int) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerName", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// ServerName indicates an expected call of ServerName.
func (mr *MockServerClientMockRecorder) ServerName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerName", reflect.TypeOf((*MockServerClient)(nil).ServerName), arg0)
}

// Shutdown mocks base method.
func (m *MockServerClient) Shutdown(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockServerClientMockRecorder) Shutdown(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockServerClient)(nil).Shutdown), arg0, arg1)
}

// Update mocks base method.
func (m *MockServerClient) Update(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockServerClientMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockServerClient)(nil).Update), arg0, arg1, arg2)
}
