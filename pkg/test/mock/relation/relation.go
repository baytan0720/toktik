// Code generated by MockGen. DO NOT EDIT.
// Source: internal/relation/kitex_gen/relation/relationservice/client.go

// Package mock_relationservice is a generated GoMock package.
package mock_relationservice

import (
	context "context"
	reflect "reflect"
	relation "toktik/internal/relation/kitex_gen/relation"

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

// Follow mocks base method.
func (m *MockClient) Follow(ctx context.Context, Req *relation.FollowReq, callOptions ...callopt.Option) (*relation.FollowRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Follow", varargs...)
	ret0, _ := ret[0].(*relation.FollowRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Follow indicates an expected call of Follow.
func (mr *MockClientMockRecorder) Follow(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Follow", reflect.TypeOf((*MockClient)(nil).Follow), varargs...)
}

// GetFollowInfo mocks base method.
func (m *MockClient) GetFollowInfo(ctx context.Context, Req *relation.GetFollowInfoReq, callOptions ...callopt.Option) (*relation.GetFollowInfoRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFollowInfo", varargs...)
	ret0, _ := ret[0].(*relation.GetFollowInfoRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFollowInfo indicates an expected call of GetFollowInfo.
func (mr *MockClientMockRecorder) GetFollowInfo(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFollowInfo", reflect.TypeOf((*MockClient)(nil).GetFollowInfo), varargs...)
}

// ListFollow mocks base method.
func (m *MockClient) ListFollow(ctx context.Context, Req *relation.ListFollowReq, callOptions ...callopt.Option) (*relation.ListFollowRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFollow", varargs...)
	ret0, _ := ret[0].(*relation.ListFollowRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFollow indicates an expected call of ListFollow.
func (mr *MockClientMockRecorder) ListFollow(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFollow", reflect.TypeOf((*MockClient)(nil).ListFollow), varargs...)
}

// ListFollower mocks base method.
func (m *MockClient) ListFollower(ctx context.Context, Req *relation.ListFollowerReq, callOptions ...callopt.Option) (*relation.ListFollowerRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFollower", varargs...)
	ret0, _ := ret[0].(*relation.ListFollowerRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFollower indicates an expected call of ListFollower.
func (mr *MockClientMockRecorder) ListFollower(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFollower", reflect.TypeOf((*MockClient)(nil).ListFollower), varargs...)
}

// ListFriend mocks base method.
func (m *MockClient) ListFriend(ctx context.Context, Req *relation.ListFriendReq, callOptions ...callopt.Option) (*relation.ListFriendRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFriend", varargs...)
	ret0, _ := ret[0].(*relation.ListFriendRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFriend indicates an expected call of ListFriend.
func (mr *MockClientMockRecorder) ListFriend(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFriend", reflect.TypeOf((*MockClient)(nil).ListFriend), varargs...)
}

// Unfollow mocks base method.
func (m *MockClient) Unfollow(ctx context.Context, Req *relation.UnfollowReq, callOptions ...callopt.Option) (*relation.UnfollowRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Unfollow", varargs...)
	ret0, _ := ret[0].(*relation.UnfollowRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unfollow indicates an expected call of Unfollow.
func (mr *MockClientMockRecorder) Unfollow(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfollow", reflect.TypeOf((*MockClient)(nil).Unfollow), varargs...)
}
