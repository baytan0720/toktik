package main

import (
	"context"
	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/relation/kitex_gen/relation"
	"toktik/internal/user/kitex_gen/user"
	"toktik/internal/user/pkg/ctx"
	"toktik/internal/video/kitex_gen/video"
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
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return nil, err
	}
	// 注册成功
	resp.UserId = userInfo.Id
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginRes, err error) {
	resp = &user.LoginRes{}

	// 验证用户名密码
	userInfo, err := s.svcCtx.UserService.Login(req.Username, req.Password)
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return nil, err
	}
	// 登录成功
	resp.UserId = userInfo.Id
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoRes, err error) {
	resp = &user.GetUserInfoRes{}
	toUserId := []int64{req.ToUserId}
	videoIdList := []int64{req.ToUserId}
	// 判断用户是否存在
	userInfo, err := s.svcCtx.UserService.GetUserInfo(req.UserId)
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return nil, err
	}
	followInfo, err := s.svcCtx.RelationClient.GetFollowInfo(ctx, &relation.GetFollowInfoReq{
		UserId:       userInfo.Id,
		ToUserIdList: toUserId,
	})
	favoriteInfo, err := s.svcCtx.FavoriteClient.GetFavoriteInfo(ctx, &favorite.GetFavoriteInfoReq{
		UserId:      userInfo.Id,
		VideoIdList: videoIdList,
	})
	videoInfo, err := s.svcCtx.VideoClient.GetVideo(ctx, &video.GetVideoReq{
		UserId:  userInfo.Id,
		VideoId: videoIdList,
	})

	totalFavorited := 0
	for i := 0; i < len(videoInfo.VideoList); i++ {
		totalFavorited = totalFavorited + int(videoInfo.VideoList[i].FavoriteCount)
	}

	favoriteCount := len(favoriteInfo.FavoriteInfoList)

	// 返回用户信息
	resp = &user.GetUserInfoRes{
		Status: user.Status_OK,
		User: &user.UserInfo{
			Id:              userInfo.Id,
			Name:            userInfo.Username,
			FollowCount:     followInfo.FollowInfoList[0].FollowCount,
			FollowerCount:   followInfo.FollowInfoList[0].FollowerCount,
			IsFollow:        followInfo.FollowInfoList[0].IsFollow,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,
			TotalFavorited:  int64(totalFavorited),
			WorkCount:       int64(len(videoInfo.VideoList)),
			FavoriteCount:   int64(favoriteCount),
		},
	}

	return resp, nil
}
