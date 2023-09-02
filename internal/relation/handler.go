package main

import (
	"context"
	"log"
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

// GetFollowInfo implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowInfo(ctx context.Context, req *relation.GetFollowInfoReq) (resp *relation.GetFollowInfoRes, err error) {
	resp = &relation.GetFollowInfoRes{}

	relations, err := s.svcCtx.RelationService.GetFollowRelations(req.UserId, req.ToUserIdList)
	if err != nil {
		return nil, err
	}

	resp.FollowInfoList = make([]*relation.FollowInfo, len(req.ToUserIdList))

	for i, j := 0, 0; i < len(req.ToUserIdList); i++ {

		followCount, err := s.svcCtx.RelationService.GetFollowCount(req.ToUserIdList[i])
		if err != nil {
			return nil, err
		}
		followerCount, err := s.svcCtx.RelationService.GetFollowerCount(req.ToUserIdList[i])
		if err != nil {
			return nil, err
		}

		if j < len(relations) && req.ToUserIdList[i] == relations[j].ToUserId {
			resp.FollowInfoList[i] = &relation.FollowInfo{
				UserId:        req.ToUserIdList[i],
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      relations[j].IsFollow,
			}
			j++
		} else {
			resp.FollowInfoList[i] = &relation.FollowInfo{
				UserId:        req.ToUserIdList[i],
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      false,
			}
		}
	}

	return
}

// Follow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Follow(ctx context.Context, req *relation.FollowReq) (resp *relation.FollowRes, err error) {
	resp = &relation.FollowRes{}

	// 判断user是否存在
	if _, err = s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
	}); err == nil {
		err = s.svcCtx.RelationService.Follow(req.UserId, req.ToUserId)
	} else {
		log.Println("get user info failed:", err)
	}
	return
}

// Unfollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Unfollow(ctx context.Context, req *relation.UnfollowReq) (resp *relation.UnfollowRes, err error) {
	resp = &relation.UnfollowRes{}

	// 判断user是否存在
	if _, err = s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
	}); err == nil {
		err = s.svcCtx.RelationService.Unfollow(req.UserId, req.ToUserId)
	} else {
		log.Println("get user info failed:", err)
	}
	return
}

// ListFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollow(ctx context.Context, req *relation.ListFollowReq) (resp *relation.ListFollowRes, err error) {
	resp = &relation.ListFollowRes{}

	userIdList, err := s.svcCtx.RelationService.GetFollow(req.UserId)
	resp.Users = make([]*relation.UserInfo, len(userIdList))
	// 遍历获取用户信息
	for i, userId := range userIdList {
		if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
			ToUserId: userId,
		}); err == nil {
			resp.Users[i] = convert2RelationUserInfo(res.User)
		} else {
			log.Println("get user info failed:", err)
		}
	}

	return
}

// ListFollower implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollower(ctx context.Context, req *relation.ListFollowerReq) (resp *relation.ListFollowerRes, err error) {
	resp = &relation.ListFollowerRes{}

	userIdList, err := s.svcCtx.RelationService.GetFollower(req.UserId)
	resp.Users = make([]*relation.UserInfo, len(userIdList))
	// 遍历获取用户信息
	for i, userId := range userIdList {
		if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
			ToUserId: userId,
		}); err == nil {
			resp.Users[i] = convert2RelationUserInfo(res.User)
		} else {
			log.Println("get user info failed:", err)
		}
	}

	return
}

// ListFriend implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFriend(ctx context.Context, req *relation.ListFriendReq) (resp *relation.ListFriendRes, err error) {
	resp = &relation.ListFriendRes{}
	userId := req.UserId

	followList, err := s.svcCtx.RelationService.GetFollow(req.UserId)
	if err != nil {
		return nil, err
	}
	followerList, err := s.svcCtx.RelationService.GetFollower(req.UserId)
	if err != nil {
		return nil, err
	}
	// 找出followList和followerList的交集
	friendList := make([]int64, 0)
	for _, follow := range followList {
		for _, follower := range followerList {
			if follow == follower {
				friendList = append(friendList, follow)
			}
		}
	}

	var messages []*message.LastMessage
	hasMessages := false

	// 获取 last message
	if res, err := s.svcCtx.MessageClient.GetLastMessage(ctx, &message.GetLastMessageReq{
		UserId:   userId,
		ToUserId: friendList,
	}); err == nil {
		messages = res.Messages
		hasMessages = true
	} else {
		log.Println("get messages failed:", err)
	}

	resp.Users = make([]*relation.FriendUser, len(friendList))
	// 遍历获取用户信息
	for i, toUserId := range friendList {
		if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
			ToUserId: toUserId,
		}); err == nil {
			resp.Users[i] = &relation.FriendUser{
				User: convert2RelationUserInfo(res.User),
			}
			if hasMessages {
				resp.Users[i].Message = messages[i].LastMessage
				resp.Users[i].MsgType = messages[i].MessageType
			}
		} else {
			log.Println("get user info failed:", err)
		}
	}

	return
}

func convert2RelationUserInfo(user *user.UserInfo) *relation.UserInfo {
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
