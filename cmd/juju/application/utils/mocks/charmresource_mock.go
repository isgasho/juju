// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/juju/application/utils (interfaces: CharmClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	charms "github.com/juju/juju/api/common/charms"
)

// MockCharmClient is a mock of CharmClient interface
type MockCharmClient struct {
	ctrl     *gomock.Controller
	recorder *MockCharmClientMockRecorder
}

// MockCharmClientMockRecorder is the mock recorder for MockCharmClient
type MockCharmClientMockRecorder struct {
	mock *MockCharmClient
}

// NewMockCharmClient creates a new mock instance
func NewMockCharmClient(ctrl *gomock.Controller) *MockCharmClient {
	mock := &MockCharmClient{ctrl: ctrl}
	mock.recorder = &MockCharmClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharmClient) EXPECT() *MockCharmClientMockRecorder {
	return m.recorder
}

// CharmInfo mocks base method
func (m *MockCharmClient) CharmInfo(arg0 string) (*charms.CharmInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmInfo", arg0)
	ret0, _ := ret[0].(*charms.CharmInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CharmInfo indicates an expected call of CharmInfo
func (mr *MockCharmClientMockRecorder) CharmInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmInfo", reflect.TypeOf((*MockCharmClient)(nil).CharmInfo), arg0)
}