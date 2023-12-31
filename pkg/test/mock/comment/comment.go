// Code generated by MockGen. DO NOT EDIT.
// Source: /home/cksfafwefasdf/gitrepo/toktik/internal/comment/kitex_gen/comment/commentservice/client.go

// Package mock_commentservice is a generated GoMock package.
package mock_commentservice

import (
	context "context"
	reflect "reflect"
	comment "toktik/internal/comment/kitex_gen/comment"

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

// CreateComment mocks base method.
func (m *MockClient) CreateComment(ctx context.Context, Req *comment.CreateCommentReq, callOptions ...callopt.Option) (*comment.CreateCommentRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateComment", varargs...)
	ret0, _ := ret[0].(*comment.CreateCommentRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment.
func (mr *MockClientMockRecorder) CreateComment(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockClient)(nil).CreateComment), varargs...)
}

// DeleteComment mocks base method.
func (m *MockClient) DeleteComment(ctx context.Context, Req *comment.DeleteCommentReq, callOptions ...callopt.Option) (*comment.DeleteCommentRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteComment", varargs...)
	ret0, _ := ret[0].(*comment.DeleteCommentRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteComment indicates an expected call of DeleteComment.
func (mr *MockClientMockRecorder) DeleteComment(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockClient)(nil).DeleteComment), varargs...)
}

// GetCommentCount mocks base method.
func (m *MockClient) GetCommentCount(ctx context.Context, Req *comment.GetCommentCountReq, callOptions ...callopt.Option) (*comment.GetCommentCountRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCommentCount", varargs...)
	ret0, _ := ret[0].(*comment.GetCommentCountRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentCount indicates an expected call of GetCommentCount.
func (mr *MockClientMockRecorder) GetCommentCount(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentCount", reflect.TypeOf((*MockClient)(nil).GetCommentCount), varargs...)
}

// ListComment mocks base method.
func (m *MockClient) ListComment(ctx context.Context, Req *comment.ListCommentReq, callOptions ...callopt.Option) (*comment.ListCommentRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListComment", varargs...)
	ret0, _ := ret[0].(*comment.ListCommentRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListComment indicates an expected call of ListComment.
func (mr *MockClientMockRecorder) ListComment(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListComment", reflect.TypeOf((*MockClient)(nil).ListComment), varargs...)
}
