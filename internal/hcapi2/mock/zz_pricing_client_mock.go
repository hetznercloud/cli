// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/cli/internal/hcapi2 (interfaces: PricingClient)

// Package hcapi2_mock is a generated GoMock package.
package hcapi2_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// MockPricingClient is a mock of PricingClient interface.
type MockPricingClient struct {
	ctrl     *gomock.Controller
	recorder *MockPricingClientMockRecorder
}

// MockPricingClientMockRecorder is the mock recorder for MockPricingClient.
type MockPricingClientMockRecorder struct {
	mock *MockPricingClient
}

// NewMockPricingClient creates a new mock instance.
func NewMockPricingClient(ctrl *gomock.Controller) *MockPricingClient {
	mock := &MockPricingClient{ctrl: ctrl}
	mock.recorder = &MockPricingClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPricingClient) EXPECT() *MockPricingClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPricingClient) Get(arg0 context.Context) (hcloud.Pricing, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(hcloud.Pricing)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockPricingClientMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPricingClient)(nil).Get), arg0)
}
