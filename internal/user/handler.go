package main

import (
	"context"
	"log"
	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/relation/kitex_gen/relation"
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
	resp = &user.RegisterRes{}

	// 判断用户名是否存在
	userInfo, err := s.svcCtx.UserService.Register(req.Username, req.Password)
	if err == nil {
		// 注册成功
		resp = &user.RegisterRes{
			UserId: userInfo.Id,
			Status: resp.Status,
		}
	} else {
		log.Println("register failed:", err)
	}
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginRes, err error) {
	resp = &user.LoginRes{}

	// 验证用户名密码
	userInfo, err := s.svcCtx.UserService.Login(req.Username, req.Password)
	if err == nil {
		// 登录成功
		resp = &user.LoginRes{
			UserId: userInfo.Id,
			Status: resp.Status,
		}
	} else {
		log.Println("login failed:", err)
	}
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoRes, err error) {
	resp = &user.GetUserInfoRes{}

	// 判断用户是否存在
	userInfo, err := s.svcCtx.UserService.GetUserInfo(req.UserId)
	followInfo, err := s.svcCtx.RelationClient.GetFollow(ctx, &relation.GetFollowInfoReq{
		UserId: userInfo.Id,
	})
	favoriteInfo, err := s.svcCtx.FavoriteClient.GetFavoriteInfo(ctx, &favorite.GetFavoriteInfoReq{
		UserId: userInfo.Id,
	})
	if err == nil {
		// 获取用户信息
		resp.User = &user.UserInfo{
			Name: userInfo.Username,

			//FollowerCount: followInfo.FollowInfoList,
			//FollowCount:   followInfo.FollowInfoList,
			//FavoriteCount: favoriteInfo.FavoriteInfoList,
			//WorkCount:
		}
	} else {
		log.Println("get user info failed:", err)
	}
	return resp, nil
}

func convert2UserInfo(user *user.UserInfo) *relation.UserInfo {
	return &relation.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        false,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
