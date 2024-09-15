// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=./service_mock.go -package=auth
//

// Package auth is a generated GoMock package.
package auth

import (
	reflect "reflect"

	domain "github.com/pietro-putelli/feynman-backend/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateUserIfNotExists mocks base method.
func (m *MockService) CreateUserIfNotExists(user *domain.ThirdPartyUser) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserIfNotExists", user)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserIfNotExists indicates an expected call of CreateUserIfNotExists.
func (mr *MockServiceMockRecorder) CreateUserIfNotExists(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserIfNotExists", reflect.TypeOf((*MockService)(nil).CreateUserIfNotExists), user)
}

// GenerateToken mocks base method.
func (m *MockService) GenerateToken(user *domain.User) (*domain.AuthTokenDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", user)
	ret0, _ := ret[0].(*domain.AuthTokenDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockServiceMockRecorder) GenerateToken(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockService)(nil).GenerateToken), user)
}

// RefreshAccessToken mocks base method.
func (m *MockService) RefreshAccessToken(refreshToken string) (*domain.AuthTokenDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshAccessToken", refreshToken)
	ret0, _ := ret[0].(*domain.AuthTokenDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshAccessToken indicates an expected call of RefreshAccessToken.
func (mr *MockServiceMockRecorder) RefreshAccessToken(refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshAccessToken", reflect.TypeOf((*MockService)(nil).RefreshAccessToken), refreshToken)
}