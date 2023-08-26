package mock_userservice

import (
	context "context"
	reflect "reflect"
	user "toktik/internal/user/kitex_gen/user"

	callopt "github.com/cloudwego/kitex/client/callopt"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetUserInfo mocks base method.
func (m *MockClient) GetUserInfo(ctx context.Context, Req *user.GetUserInfoReq, callOptions ...callopt.Option) (*user.GetUserInfoRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserInfo", varargs...)
	ret0, _ := ret[0].(*user.GetUserInfoRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfo indicates an expected call of GetUserInfo.
func (mr *MockClientMockRecorder) GetUserInfo(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*MockClient)(nil).GetUserInfo), varargs...)
}

// Login mocks base method.
func (m *MockClient) Login(ctx context.Context, Req *user.LoginReq, callOptions ...callopt.Option) (*user.LoginRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Login", varargs...)
	ret0, _ := ret[0].(*user.LoginRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockClientMockRecorder) Login(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockClient)(nil).Login), varargs...)
}

// Register mocks base method.
func (m *MockClient) Register(ctx context.Context, Req *user.RegisterReq, callOptions ...callopt.Option) (*user.RegisterRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Register", varargs...)
	ret0, _ := ret[0].(*user.RegisterRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockClientMockRecorder) Register(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockClient)(nil).Register), varargs...)
}
