package main

import (
	"context"
	"sync"

	"toktik/internal/message/kitex_gen/message"
	"toktik/internal/relation/kitex_gen/relation"
	"toktik/internal/relation/pkg/ctx"
	"toktik/internal/user/kitex_gen/user"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewRelationServiceImpl(svcCtx *ctx.ServiceContext) *RelationServiceImpl {
	return &RelationServiceImpl{
		svcCtx: svcCtx,
	}
}

// Error messages.
var (
	errGetFollowInfoFailed    = "获取关注信息失败"
	errGetFollowCountFailed   = "获取关注数量失败"
	errGetFollowerCountFailed = "获取粉丝数量失败"
	errFollowFailed           = "关注失败"
	errUnfollowFailed         = "取消关注失败"
	errListFollowFailed       = "获取关注列表失败"
	errListFollowerFailed     = "获取粉丝列表失败"
	errIsFollowFailed         = "判断是否关注失败"
	errGetUserInfoFailed      = "获取用户信息失败"
)

// GetFollowInfo implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowInfo(ctx context.Context, req *relation.GetFollowInfoReq) (resp *relation.GetFollowInfoRes, _ error) {
	resp = &relation.GetFollowInfoRes{}

	isFollowMap, err := s.svcCtx.RelationService.GetFollows(req.UserId, req.ToUserIdList)
	if err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errGetFollowInfoFailed
		return
	}

	resp.FollowInfoList = make([]*relation.FollowInfo, 0, len(req.ToUserIdList))

	for _, toUserId := range req.ToUserIdList {
		followCount, err := s.svcCtx.RelationService.GetFollowCount(toUserId)
		if err != nil {
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = errGetFollowCountFailed
			return
		}
		followerCount, err := s.svcCtx.RelationService.GetFollowerCount(toUserId)
		if err != nil {
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = errGetFollowerCountFailed
			return
		}
		resp.FollowInfoList = append(resp.FollowInfoList, &relation.FollowInfo{
			UserId:        req.UserId,
			IsFollow:      isFollowMap[toUserId],
			FollowCount:   followCount,
			FollowerCount: followerCount,
		})
	}

	return
}

// Follow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Follow(ctx context.Context, req *relation.FollowReq) (resp *relation.FollowRes, _ error) {
	resp = &relation.FollowRes{}

	// 获取用户信息
	if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
	}); err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errGetUserInfoFailed
		return
	} else if res.Status != user.Status_OK {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = res.ErrMsg
		return
	}
	err := s.svcCtx.RelationService.Follow(req.UserId, req.ToUserId)
	if err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errFollowFailed
	}
	return
}

// Unfollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Unfollow(ctx context.Context, req *relation.UnfollowReq) (resp *relation.UnfollowRes, _ error) {
	resp = &relation.UnfollowRes{}
	err := s.svcCtx.RelationService.Unfollow(req.UserId, req.ToUserId)
	if err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errUnfollowFailed
	}
	return
}

// ListFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollow(ctx context.Context, req *relation.ListFollowReq) (resp *relation.ListFollowRes, _ error) {
	resp = &relation.ListFollowRes{}

	userIdList, err := s.svcCtx.RelationService.ListFollow(req.UserId)
	if err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errListFollowFailed
		return
	}

	resp.Users = make([]*relation.UserInfo, 0, len(userIdList))
	// 获取用户信息
	userId2UserInfo := make(map[int64]*relation.UserInfo)
	for _, toUserId := range userIdList {
		followCount, err := s.svcCtx.RelationService.GetFollowCount(toUserId)
		if err != nil {
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = errGetFollowCountFailed
			return
		}
		followerCount, err := s.svcCtx.RelationService.GetFollowerCount(toUserId)
		if err != nil {
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = errGetFollowerCountFailed
			return
		}
		userId2UserInfo[toUserId] = &relation.UserInfo{
			Id:            toUserId,
			IsFollow:      true,
			FollowCount:   followCount,
			FollowerCount: followerCount,
		}
		resp.Users = append(resp.Users, userId2UserInfo[toUserId])
	}
	if res, err := s.svcCtx.UserClient.GetUserInfos(ctx, &user.GetUserInfosReq{
		ToUserIds: userIdList,
	}); err == nil && res.Status == user.Status_OK {
		for i, userInfo := range res.Users {
			userInfo.IsFollow = true
			userInfo.FollowCount = userId2UserInfo[userInfo.Id].FollowCount
			userInfo.FollowerCount = userId2UserInfo[userInfo.Id].FollowerCount
			resp.Users[i] = convert2RelationUserInfo(userInfo)
		}
	}
	return
}

// ListFollower implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollower(ctx context.Context, req *relation.ListFollowerReq) (resp *relation.ListFollowerRes, _ error) {
	resp = &relation.ListFollowerRes{}

	userIdList, err := s.svcCtx.RelationService.ListFollower(req.UserId)
	if err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errListFollowerFailed
		return
	}

	resp.Users = make([]*relation.UserInfo, 0, len(userIdList))
	// 获取用户信息
	userId2UserInfo := make(map[int64]*relation.UserInfo)
	for _, toUserId := range userIdList {
		isFollow, err := s.svcCtx.RelationService.IsFollow(req.UserId, toUserId)
		if err != nil {
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = errIsFollowFailed
			return
		}
		followCount, err := s.svcCtx.RelationService.GetFollowCount(toUserId)
		if err != nil {
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = errGetFollowCountFailed
			return
		}
		followerCount, err := s.svcCtx.RelationService.GetFollowerCount(toUserId)
		if err != nil {
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = errGetFollowerCountFailed
			return
		}
		userId2UserInfo[toUserId] = &relation.UserInfo{
			Id:            toUserId,
			IsFollow:      isFollow,
			FollowCount:   followCount,
			FollowerCount: followerCount,
		}
		resp.Users = append(resp.Users, userId2UserInfo[toUserId])
	}
	if res, err := s.svcCtx.UserClient.GetUserInfos(ctx, &user.GetUserInfosReq{
		ToUserIds: userIdList,
	}); err == nil && res.Status == user.Status_OK {
		for i, userInfo := range res.Users {
			userInfo.IsFollow = userId2UserInfo[userInfo.Id].IsFollow
			userInfo.FollowCount = userId2UserInfo[userInfo.Id].FollowCount
			userInfo.FollowerCount = userId2UserInfo[userInfo.Id].FollowerCount
			resp.Users[i] = convert2RelationUserInfo(userInfo)
		}
	}
	return
}

// ListFriend implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFriend(ctx context.Context, req *relation.ListFriendReq) (resp *relation.ListFriendRes, _ error) {
	resp = &relation.ListFriendRes{}
	userId := req.UserId

	followList, err := s.svcCtx.RelationService.ListFollow(req.UserId)
	if err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errListFollowFailed
		return
	}
	followerList, err := s.svcCtx.RelationService.ListFollower(req.UserId)
	if err != nil {
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = errListFollowerFailed
		return
	}

	followMap := make(map[int64]bool)
	friendList := make([]int64, 0, len(followList))
	// 找出followList和followerList的交集
	for _, toUserId := range followList {
		followMap[toUserId] = true
	}
	for _, toUserId := range followerList {
		if followMap[toUserId] {
			friendList = append(friendList, toUserId)
		}
	}

	wg := sync.WaitGroup{}

	// 获取 last message
	userId2LastMessage := make(map[int64]*message.LastMessage)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, toUserId := range friendList {
			userId2LastMessage[toUserId] = &message.LastMessage{
				ToUserId: toUserId,
			}
		}

		if res, err := s.svcCtx.MessageClient.GetLastMessage(ctx, &message.GetLastMessageReq{
			UserId:   userId,
			ToUserId: friendList,
		}); err == nil && res.Status == message.Status_OK {
			for _, lastMessage := range res.Messages {
				userId2LastMessage[lastMessage.ToUserId] = lastMessage
			}
		}
	}()

	// 获取用户信息
	userId2UserInfo := make(map[int64]*relation.UserInfo)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, toUserId := range friendList {
			followCount, err := s.svcCtx.RelationService.GetFollowCount(userId)
			if err != nil {
				resp.Status = relation.Status_ERROR
				resp.ErrMsg = errGetFollowCountFailed
				return
			}
			followerCount, err := s.svcCtx.RelationService.GetFollowerCount(userId)
			if err != nil {
				resp.Status = relation.Status_ERROR
				resp.ErrMsg = errGetFollowerCountFailed
				return
			}
			userId2UserInfo[toUserId] = &relation.UserInfo{
				Id:            toUserId,
				IsFollow:      true,
				FollowCount:   followCount,
				FollowerCount: followerCount,
			}
		}

		if res, err := s.svcCtx.UserClient.GetUserInfos(ctx, &user.GetUserInfosReq{
			ToUserIds: friendList,
		}); err == nil && res.Status == user.Status_OK {
			for _, userInfo := range res.Users {
				userInfo.IsFollow = true
				userInfo.FollowCount = userId2UserInfo[userInfo.Id].FollowCount
				userInfo.FollowerCount = userId2UserInfo[userInfo.Id].FollowerCount
				userId2UserInfo[userInfo.Id] = convert2RelationUserInfo(userInfo)
			}
		}
	}()

	wg.Wait()

	resp.Users = make([]*relation.FriendUser, 0, len(friendList))
	for _, toUserId := range friendList {
		friendUser := &relation.FriendUser{
			User:    userId2UserInfo[toUserId],
			Message: userId2LastMessage[toUserId].LastMessage,
			MsgType: userId2LastMessage[toUserId].MessageType,
		}
		resp.Users = append(resp.Users, friendUser)
	}

	return
}

func convert2RelationUserInfo(user *user.UserInfo) *relation.UserInfo {
	return &relation.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
