package main

import (
	"context"
	"sync"

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
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterRes, _ error) {
	resp = &user.RegisterRes{}
	userId, err := s.svcCtx.UserService.CreateUser(req.Username, req.Password)
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}
	resp.UserId = userId
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginRes, _ error) {
	resp = &user.LoginRes{}
	userId, err := s.svcCtx.UserService.Login(req.Username, req.Password)
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}
	resp.UserId = userId
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoRes, _ error) {
	resp = &user.GetUserInfoRes{}
	userInfo, err := s.svcCtx.UserService.GetUserById(req.ToUserId)
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	resp.User = &user.UserInfo{
		Id:              userInfo.Id,
		Name:            userInfo.Username,
		Avatar:          userInfo.Avatar,
		BackgroundImage: userInfo.BackgroundImage,
		Signature:       userInfo.Signature,
	}

	wg := &sync.WaitGroup{}

	// get user relation info
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.RelationClient.GetFollowInfo(ctx, &relation.GetFollowInfoReq{
			UserId:       req.UserId,
			ToUserIdList: []int64{req.ToUserId},
		})
		if err != nil || res.Status != relation.Status_OK {
			return
		}
		followInfo := res.FollowInfoList[0]
		resp.User.FollowCount = followInfo.FollowCount
		resp.User.FollowerCount = followInfo.FollowerCount
		resp.User.IsFollow = followInfo.IsFollow
	}()

	// get user favorite info
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.FavoriteClient.GetUserFavoriteInfo(ctx, &favorite.GetUserFavoriteInfoReq{
			UserIdList: []int64{req.ToUserId},
		})
		if err != nil || res.Status != favorite.Status_OK {
			return
		}
		favoriteInfo := res.FavoriteInfoList[0]
		resp.User.FavoriteCount = favoriteInfo.FavoriteCount
		resp.User.TotalFavorited = favoriteInfo.TotalFavorited
	}()

	// get user work count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.VideoClient.GetWorkCount(ctx, &video.GetWorkCountReq{
			UserIdList: []int64{req.ToUserId},
		})
		if err != nil || res.Status != video.Status_OK {
			return
		}
		workCountInfo := res.WorkCountList[0]
		resp.User.WorkCount = workCountInfo.WorkCount
	}()

	wg.Wait()

	return
}

// GetUserInfos implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfos(ctx context.Context, req *user.GetUserInfosReq) (resp *user.GetUserInfosRes, _ error) {
	resp = &user.GetUserInfosRes{}
	userInfos, err := s.svcCtx.UserService.GetUserByIds(req.ToUserIds)
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	userId2UserInfo := make(map[int64]*user.UserInfo)
	resp.Users = make([]*user.UserInfo, 0, len(userInfos))
	for _, userInfo := range userInfos {
		user := &user.UserInfo{
			Id:              userInfo.Id,
			Name:            userInfo.Username,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,
		}
		resp.Users = append(resp.Users, user)
		userId2UserInfo[userInfo.Id] = user
	}

	// get user relation info
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()

		res, err := s.svcCtx.RelationClient.GetFollowInfo(ctx, &relation.GetFollowInfoReq{
			UserId:       req.UserId,
			ToUserIdList: req.ToUserIds,
		})
		if err != nil || res.Status != relation.Status_OK {
			return
		}
		for _, followInfo := range res.FollowInfoList {
			userInfo := userId2UserInfo[followInfo.UserId]
			userInfo.FollowCount = followInfo.FollowCount
			userInfo.FollowerCount = followInfo.FollowerCount
			userInfo.IsFollow = followInfo.IsFollow
		}
	}()

	// get user favorite info
	go func() {
		wg.Add(1)
		defer wg.Done()

		res, err := s.svcCtx.FavoriteClient.GetUserFavoriteInfo(ctx, &favorite.GetUserFavoriteInfoReq{
			UserIdList: req.ToUserIds,
		})
		if err != nil || res.Status != favorite.Status_OK {
			return
		}
		for _, favoriteInfo := range res.FavoriteInfoList {
			userInfo := userId2UserInfo[favoriteInfo.UserId]
			userInfo.FavoriteCount = favoriteInfo.FavoriteCount
			userInfo.TotalFavorited = favoriteInfo.TotalFavorited
		}
	}()

	// get user work count
	go func() {
		wg.Add(1)
		defer wg.Done()

		res, err := s.svcCtx.VideoClient.GetWorkCount(ctx, &video.GetWorkCountReq{
			UserIdList: []int64{req.UserId},
		})
		if err != nil || res.Status != video.Status_OK {
			return
		}
		for _, workCountInfo := range res.WorkCountList {
			userInfo := userId2UserInfo[workCountInfo.UserId]
			userInfo.WorkCount = workCountInfo.WorkCount
		}
	}()

	wg.Wait()

	return
}
