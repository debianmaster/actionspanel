// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/phunki/actionspanel/pkg/gh (interfaces: AppsService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v30/github"
	reflect "reflect"
)

// MockAppsService is a mock of AppsService interface
type MockAppsService struct {
	ctrl     *gomock.Controller
	recorder *MockAppsServiceMockRecorder
}

// MockAppsServiceMockRecorder is the mock recorder for MockAppsService
type MockAppsServiceMockRecorder struct {
	mock *MockAppsService
}

// NewMockAppsService creates a new mock instance
func NewMockAppsService(ctrl *gomock.Controller) *MockAppsService {
	mock := &MockAppsService{ctrl: ctrl}
	mock.recorder = &MockAppsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppsService) EXPECT() *MockAppsServiceMockRecorder {
	return m.recorder
}

// ListUserInstallations mocks base method
func (m *MockAppsService) ListUserInstallations(arg0 context.Context, arg1 *github.ListOptions) ([]*github.Installation, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserInstallations", arg0, arg1)
	ret0, _ := ret[0].([]*github.Installation)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListUserInstallations indicates an expected call of ListUserInstallations
func (mr *MockAppsServiceMockRecorder) ListUserInstallations(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserInstallations", reflect.TypeOf((*MockAppsService)(nil).ListUserInstallations), arg0, arg1)
}
