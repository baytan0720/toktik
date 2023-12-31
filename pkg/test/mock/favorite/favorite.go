// Code generated by MockGen. DO NOT EDIT.
// Source: internal/favorite/kitex_gen/favorite/favoriteservice/client.go

// Package mock_favoriteservice is a generated GoMock package.
package mock_favoriteservice

import (
	context "context"
	reflect "reflect"
	favorite "toktik/internal/favorite/kitex_gen/favorite"

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

// Favorite mocks base method.
func (m *MockClient) Favorite(ctx context.Context, Req *favorite.FavoriteReq, callOptions ...callopt.Option) (*favorite.FavoriteRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Favorite", varargs...)
	ret0, _ := ret[0].(*favorite.FavoriteRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Favorite indicates an expected call of Favorite.
func (mr *MockClientMockRecorder) Favorite(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Favorite", reflect.TypeOf((*MockClient)(nil).Favorite), varargs...)
}

// GetUserFavoriteInfo mocks base method.
func (m *MockClient) GetUserFavoriteInfo(ctx context.Context, Req *favorite.GetUserFavoriteInfoReq, callOptions ...callopt.Option) (*favorite.GetUserFavoriteInfoRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserFavoriteInfo", varargs...)
	ret0, _ := ret[0].(*favorite.GetUserFavoriteInfoRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFavoriteInfo indicates an expected call of GetUserFavoriteInfo.
func (mr *MockClientMockRecorder) GetUserFavoriteInfo(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFavoriteInfo", reflect.TypeOf((*MockClient)(nil).GetUserFavoriteInfo), varargs...)
}

// GetVideoFavoriteInfo mocks base method.
func (m *MockClient) GetVideoFavoriteInfo(ctx context.Context, Req *favorite.GetVideoFavoriteInfoReq, callOptions ...callopt.Option) (*favorite.GetVideoFavoriteInfoRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetVideoFavoriteInfo", varargs...)
	ret0, _ := ret[0].(*favorite.GetVideoFavoriteInfoRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideoFavoriteInfo indicates an expected call of GetVideoFavoriteInfo.
func (mr *MockClientMockRecorder) GetVideoFavoriteInfo(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideoFavoriteInfo", reflect.TypeOf((*MockClient)(nil).GetVideoFavoriteInfo), varargs...)
}

// ListFavorite mocks base method.
func (m *MockClient) ListFavorite(ctx context.Context, Req *favorite.ListFavoriteReq, callOptions ...callopt.Option) (*favorite.ListFavoriteRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFavorite", varargs...)
	ret0, _ := ret[0].(*favorite.ListFavoriteRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFavorite indicates an expected call of ListFavorite.
func (mr *MockClientMockRecorder) ListFavorite(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFavorite", reflect.TypeOf((*MockClient)(nil).ListFavorite), varargs...)
}

// UnFavorite mocks base method.
func (m *MockClient) UnFavorite(ctx context.Context, Req *favorite.UnFavoriteReq, callOptions ...callopt.Option) (*favorite.UnFavoriteRes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, Req}
	for _, a := range callOptions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UnFavorite", varargs...)
	ret0, _ := ret[0].(*favorite.UnFavoriteRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnFavorite indicates an expected call of UnFavorite.
func (mr *MockClientMockRecorder) UnFavorite(ctx, Req interface{}, callOptions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, Req}, callOptions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnFavorite", reflect.TypeOf((*MockClient)(nil).UnFavorite), varargs...)
}
