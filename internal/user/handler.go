package main

import (
	"context"

	"toktik/internal/user/kitex_gen/user"
	"toktik/internal/user/pkg/ctx"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewUserServiceImpl(svcCtx *ctx.ServiceContext) *UserServiceImpl {
	return &UserServiceImpl{
		svcCtx: svcCtx,
	}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterRes, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginRes, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoRes, err error) {
	// TODO: Your code here...
	return
}
