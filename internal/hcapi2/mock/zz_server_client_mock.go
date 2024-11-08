// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: ServerClient)
//
// Generated by this command:
//
//	mockgen -package mock -destination zz_server_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ServerClient
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// MockServerClient is a mock of ServerClient interface.
type MockServerClient struct {
	ctrl     *gomock.Controller
	recorder *MockServerClientMockRecorder
	isgomock struct{}
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

// AddToPlacementGroup mocks base method.
func (m *MockServerClient) AddToPlacementGroup(ctx context.Context, server *hcloud.Server, placementGroup *hcloud.PlacementGroup) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToPlacementGroup", ctx, server, placementGroup)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddToPlacementGroup indicates an expected call of AddToPlacementGroup.
func (mr *MockServerClientMockRecorder) AddToPlacementGroup(ctx, server, placementGroup any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToPlacementGroup", reflect.TypeOf((*MockServerClient)(nil).AddToPlacementGroup), ctx, server, placementGroup)
}

// All mocks base method.
func (m *MockServerClient) All(ctx context.Context) ([]*hcloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", ctx)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockServerClientMockRecorder) All(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockServerClient)(nil).All), ctx)
}

// AllWithOpts mocks base method.
func (m *MockServerClient) AllWithOpts(ctx context.Context, opts hcloud.ServerListOpts) ([]*hcloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", ctx, opts)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockServerClientMockRecorder) AllWithOpts(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockServerClient)(nil).AllWithOpts), ctx, opts)
}

// AttachISO mocks base method.
func (m *MockServerClient) AttachISO(ctx context.Context, server *hcloud.Server, iso *hcloud.ISO) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachISO", ctx, server, iso)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachISO indicates an expected call of AttachISO.
func (mr *MockServerClientMockRecorder) AttachISO(ctx, server, iso any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachISO", reflect.TypeOf((*MockServerClient)(nil).AttachISO), ctx, server, iso)
}

// AttachToNetwork mocks base method.
func (m *MockServerClient) AttachToNetwork(ctx context.Context, server *hcloud.Server, opts hcloud.ServerAttachToNetworkOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachToNetwork", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachToNetwork indicates an expected call of AttachToNetwork.
func (mr *MockServerClientMockRecorder) AttachToNetwork(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachToNetwork", reflect.TypeOf((*MockServerClient)(nil).AttachToNetwork), ctx, server, opts)
}

// ChangeAliasIPs mocks base method.
func (m *MockServerClient) ChangeAliasIPs(ctx context.Context, server *hcloud.Server, opts hcloud.ServerChangeAliasIPsOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAliasIPs", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeAliasIPs indicates an expected call of ChangeAliasIPs.
func (mr *MockServerClientMockRecorder) ChangeAliasIPs(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAliasIPs", reflect.TypeOf((*MockServerClient)(nil).ChangeAliasIPs), ctx, server, opts)
}

// ChangeDNSPtr mocks base method.
func (m *MockServerClient) ChangeDNSPtr(ctx context.Context, server *hcloud.Server, ip string, ptr *string) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDNSPtr", ctx, server, ip, ptr)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeDNSPtr indicates an expected call of ChangeDNSPtr.
func (mr *MockServerClientMockRecorder) ChangeDNSPtr(ctx, server, ip, ptr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDNSPtr", reflect.TypeOf((*MockServerClient)(nil).ChangeDNSPtr), ctx, server, ip, ptr)
}

// ChangeProtection mocks base method.
func (m *MockServerClient) ChangeProtection(ctx context.Context, server *hcloud.Server, opts hcloud.ServerChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeProtection", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeProtection indicates an expected call of ChangeProtection.
func (mr *MockServerClientMockRecorder) ChangeProtection(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeProtection", reflect.TypeOf((*MockServerClient)(nil).ChangeProtection), ctx, server, opts)
}

// ChangeType mocks base method.
func (m *MockServerClient) ChangeType(ctx context.Context, server *hcloud.Server, opts hcloud.ServerChangeTypeOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeType", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeType indicates an expected call of ChangeType.
func (mr *MockServerClientMockRecorder) ChangeType(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeType", reflect.TypeOf((*MockServerClient)(nil).ChangeType), ctx, server, opts)
}

// Create mocks base method.
func (m *MockServerClient) Create(ctx context.Context, opts hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, opts)
	ret0, _ := ret[0].(hcloud.ServerCreateResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockServerClientMockRecorder) Create(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockServerClient)(nil).Create), ctx, opts)
}

// CreateImage mocks base method.
func (m *MockServerClient) CreateImage(ctx context.Context, server *hcloud.Server, opts *hcloud.ServerCreateImageOpts) (hcloud.ServerCreateImageResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImage", ctx, server, opts)
	ret0, _ := ret[0].(hcloud.ServerCreateImageResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateImage indicates an expected call of CreateImage.
func (mr *MockServerClientMockRecorder) CreateImage(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImage", reflect.TypeOf((*MockServerClient)(nil).CreateImage), ctx, server, opts)
}

// Delete mocks base method.
func (m *MockServerClient) Delete(ctx context.Context, server *hcloud.Server) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, server)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockServerClientMockRecorder) Delete(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockServerClient)(nil).Delete), ctx, server)
}

// DeleteWithResult mocks base method.
func (m *MockServerClient) DeleteWithResult(ctx context.Context, server *hcloud.Server) (*hcloud.ServerDeleteResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWithResult", ctx, server)
	ret0, _ := ret[0].(*hcloud.ServerDeleteResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteWithResult indicates an expected call of DeleteWithResult.
func (mr *MockServerClientMockRecorder) DeleteWithResult(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWithResult", reflect.TypeOf((*MockServerClient)(nil).DeleteWithResult), ctx, server)
}

// DetachFromNetwork mocks base method.
func (m *MockServerClient) DetachFromNetwork(ctx context.Context, server *hcloud.Server, opts hcloud.ServerDetachFromNetworkOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachFromNetwork", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DetachFromNetwork indicates an expected call of DetachFromNetwork.
func (mr *MockServerClientMockRecorder) DetachFromNetwork(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachFromNetwork", reflect.TypeOf((*MockServerClient)(nil).DetachFromNetwork), ctx, server, opts)
}

// DetachISO mocks base method.
func (m *MockServerClient) DetachISO(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachISO", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DetachISO indicates an expected call of DetachISO.
func (mr *MockServerClientMockRecorder) DetachISO(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachISO", reflect.TypeOf((*MockServerClient)(nil).DetachISO), ctx, server)
}

// DisableBackup mocks base method.
func (m *MockServerClient) DisableBackup(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableBackup", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DisableBackup indicates an expected call of DisableBackup.
func (mr *MockServerClientMockRecorder) DisableBackup(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableBackup", reflect.TypeOf((*MockServerClient)(nil).DisableBackup), ctx, server)
}

// DisableRescue mocks base method.
func (m *MockServerClient) DisableRescue(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableRescue", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DisableRescue indicates an expected call of DisableRescue.
func (mr *MockServerClientMockRecorder) DisableRescue(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableRescue", reflect.TypeOf((*MockServerClient)(nil).DisableRescue), ctx, server)
}

// EnableBackup mocks base method.
func (m *MockServerClient) EnableBackup(ctx context.Context, server *hcloud.Server, window string) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableBackup", ctx, server, window)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EnableBackup indicates an expected call of EnableBackup.
func (mr *MockServerClientMockRecorder) EnableBackup(ctx, server, window any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableBackup", reflect.TypeOf((*MockServerClient)(nil).EnableBackup), ctx, server, window)
}

// EnableRescue mocks base method.
func (m *MockServerClient) EnableRescue(ctx context.Context, server *hcloud.Server, opts hcloud.ServerEnableRescueOpts) (hcloud.ServerEnableRescueResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableRescue", ctx, server, opts)
	ret0, _ := ret[0].(hcloud.ServerEnableRescueResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EnableRescue indicates an expected call of EnableRescue.
func (mr *MockServerClientMockRecorder) EnableRescue(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableRescue", reflect.TypeOf((*MockServerClient)(nil).EnableRescue), ctx, server, opts)
}

// Get mocks base method.
func (m *MockServerClient) Get(ctx context.Context, idOrName string) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, idOrName)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockServerClientMockRecorder) Get(ctx, idOrName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockServerClient)(nil).Get), ctx, idOrName)
}

// GetByID mocks base method.
func (m *MockServerClient) GetByID(ctx context.Context, id int64) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockServerClientMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockServerClient)(nil).GetByID), ctx, id)
}

// GetByName mocks base method.
func (m *MockServerClient) GetByName(ctx context.Context, name string) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockServerClientMockRecorder) GetByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockServerClient)(nil).GetByName), ctx, name)
}

// GetMetrics mocks base method.
func (m *MockServerClient) GetMetrics(ctx context.Context, server *hcloud.Server, opts hcloud.ServerGetMetricsOpts) (*hcloud.ServerMetrics, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.ServerMetrics)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockServerClientMockRecorder) GetMetrics(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockServerClient)(nil).GetMetrics), ctx, server, opts)
}

// LabelKeys mocks base method.
func (m *MockServerClient) LabelKeys(idOrName string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LabelKeys", idOrName)
	ret0, _ := ret[0].([]string)
	return ret0
}

// LabelKeys indicates an expected call of LabelKeys.
func (mr *MockServerClientMockRecorder) LabelKeys(idOrName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LabelKeys", reflect.TypeOf((*MockServerClient)(nil).LabelKeys), idOrName)
}

// List mocks base method.
func (m *MockServerClient) List(ctx context.Context, opts hcloud.ServerListOpts) ([]*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, opts)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockServerClientMockRecorder) List(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockServerClient)(nil).List), ctx, opts)
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
func (m *MockServerClient) Poweroff(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poweroff", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Poweroff indicates an expected call of Poweroff.
func (mr *MockServerClientMockRecorder) Poweroff(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poweroff", reflect.TypeOf((*MockServerClient)(nil).Poweroff), ctx, server)
}

// Poweron mocks base method.
func (m *MockServerClient) Poweron(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poweron", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Poweron indicates an expected call of Poweron.
func (mr *MockServerClientMockRecorder) Poweron(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poweron", reflect.TypeOf((*MockServerClient)(nil).Poweron), ctx, server)
}

// Reboot mocks base method.
func (m *MockServerClient) Reboot(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reboot", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reboot indicates an expected call of Reboot.
func (mr *MockServerClientMockRecorder) Reboot(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reboot", reflect.TypeOf((*MockServerClient)(nil).Reboot), ctx, server)
}

// Rebuild mocks base method.
func (m *MockServerClient) Rebuild(ctx context.Context, server *hcloud.Server, opts hcloud.ServerRebuildOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rebuild", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Rebuild indicates an expected call of Rebuild.
func (mr *MockServerClientMockRecorder) Rebuild(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rebuild", reflect.TypeOf((*MockServerClient)(nil).Rebuild), ctx, server, opts)
}

// RebuildWithResult mocks base method.
func (m *MockServerClient) RebuildWithResult(ctx context.Context, server *hcloud.Server, opts hcloud.ServerRebuildOpts) (hcloud.ServerRebuildResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RebuildWithResult", ctx, server, opts)
	ret0, _ := ret[0].(hcloud.ServerRebuildResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RebuildWithResult indicates an expected call of RebuildWithResult.
func (mr *MockServerClientMockRecorder) RebuildWithResult(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RebuildWithResult", reflect.TypeOf((*MockServerClient)(nil).RebuildWithResult), ctx, server, opts)
}

// RemoveFromPlacementGroup mocks base method.
func (m *MockServerClient) RemoveFromPlacementGroup(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromPlacementGroup", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RemoveFromPlacementGroup indicates an expected call of RemoveFromPlacementGroup.
func (mr *MockServerClientMockRecorder) RemoveFromPlacementGroup(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromPlacementGroup", reflect.TypeOf((*MockServerClient)(nil).RemoveFromPlacementGroup), ctx, server)
}

// RequestConsole mocks base method.
func (m *MockServerClient) RequestConsole(ctx context.Context, server *hcloud.Server) (hcloud.ServerRequestConsoleResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestConsole", ctx, server)
	ret0, _ := ret[0].(hcloud.ServerRequestConsoleResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RequestConsole indicates an expected call of RequestConsole.
func (mr *MockServerClientMockRecorder) RequestConsole(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestConsole", reflect.TypeOf((*MockServerClient)(nil).RequestConsole), ctx, server)
}

// Reset mocks base method.
func (m *MockServerClient) Reset(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reset indicates an expected call of Reset.
func (mr *MockServerClientMockRecorder) Reset(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockServerClient)(nil).Reset), ctx, server)
}

// ResetPassword mocks base method.
func (m *MockServerClient) ResetPassword(ctx context.Context, server *hcloud.Server) (hcloud.ServerResetPasswordResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", ctx, server)
	ret0, _ := ret[0].(hcloud.ServerResetPasswordResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockServerClientMockRecorder) ResetPassword(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockServerClient)(nil).ResetPassword), ctx, server)
}

// ServerName mocks base method.
func (m *MockServerClient) ServerName(id int64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerName", id)
	ret0, _ := ret[0].(string)
	return ret0
}

// ServerName indicates an expected call of ServerName.
func (mr *MockServerClientMockRecorder) ServerName(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerName", reflect.TypeOf((*MockServerClient)(nil).ServerName), id)
}

// Shutdown mocks base method.
func (m *MockServerClient) Shutdown(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", ctx, server)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockServerClientMockRecorder) Shutdown(ctx, server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockServerClient)(nil).Shutdown), ctx, server)
}

// Update mocks base method.
func (m *MockServerClient) Update(ctx context.Context, server *hcloud.Server, opts hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, server, opts)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockServerClientMockRecorder) Update(ctx, server, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockServerClient)(nil).Update), ctx, server, opts)
}
