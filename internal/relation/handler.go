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
	isFollowMap := make(map[int64]bool)
	for _, r := range relations {
		isFollowMap[r.ToUserId] = r.IsFollow
	}

	resp.FollowInfoList = make([]*relation.FollowInfo, 0, len(req.ToUserIdList))

	for _, toUserId := range req.ToUserIdList {
		followCount, err := s.svcCtx.RelationService.GetFollowCount(toUserId)
		if err != nil {
			return nil, err
		}
		followerCount, err := s.svcCtx.RelationService.GetFollowerCount(toUserId)
		if err != nil {
			return nil, err
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
func (s *RelationServiceImpl) Follow(ctx context.Context, req *relation.FollowReq) (resp *relation.FollowRes, err error) {
	resp = &relation.FollowRes{}

	// 判断user是否存在
	if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
	}); err != nil {
		log.Println("get user info failed:", err)
		return nil, err
	} else if res.Status != user.Status_OK {
		log.Println("user not exist")
		resp.Status = relation.Status_ERROR
		resp.ErrMsg = "user not exist"
	} else {
		err := s.svcCtx.RelationService.Follow(req.UserId, req.ToUserId)
		return resp, err
	}
	return
}

// Unfollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Unfollow(ctx context.Context, req *relation.UnfollowReq) (resp *relation.UnfollowRes, err error) {
	resp = &relation.UnfollowRes{}
	err = s.svcCtx.RelationService.Unfollow(req.UserId, req.ToUserId)
	return
}

// ListFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollow(ctx context.Context, req *relation.ListFollowReq) (resp *relation.ListFollowRes, err error) {
	resp = &relation.ListFollowRes{}

	userIdList, err := s.svcCtx.RelationService.ListFollow(req.UserId)
	if err != nil {
		return nil, err
	}

	resp.Users = make([]*relation.UserInfo, len(userIdList))
	// 遍历获取用户信息
	for i, userId := range userIdList {
		if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
			ToUserId: userId,
		}); err != nil {
			log.Println("get user info failed:", err)
		} else {
			followCount, err := s.svcCtx.RelationService.GetFollowCount(userId)
			if err != nil {
				return nil, err
			}
			followerCount, err := s.svcCtx.RelationService.GetFollowerCount(userId)
			if err != nil {
				return nil, err
			}
			resp.Users[i] = convert2RelationUserInfo(res.User)
			resp.Users[i].IsFollow = true
			resp.Users[i].FollowCount = followCount
			resp.Users[i].FollowerCount = followerCount
		}
	}

	return
}

// ListFollower implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollower(ctx context.Context, req *relation.ListFollowerReq) (resp *relation.ListFollowerRes, err error) {
	resp = &relation.ListFollowerRes{}

	userIdList, err := s.svcCtx.RelationService.ListFollower(req.UserId)
	resp.Users = make([]*relation.UserInfo, len(userIdList))
	// 遍历获取用户信息
	for i, userId := range userIdList {
		if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
			ToUserId: userId,
		}); err != nil {
			log.Println("get user info failed:", err)
		} else {
			isFollow, err := s.svcCtx.RelationService.IsFollow(req.UserId, userId)
			if err != nil {
				return nil, err
			}
			followCount, err := s.svcCtx.RelationService.GetFollowCount(userId)
			if err != nil {
				return nil, err
			}
			followerCount, err := s.svcCtx.RelationService.GetFollowerCount(userId)
			if err != nil {
				return nil, err
			}
			resp.Users[i] = convert2RelationUserInfo(res.User)
			resp.Users[i].IsFollow = isFollow
			resp.Users[i].FollowCount = followCount
			resp.Users[i].FollowerCount = followerCount
		}
	}

	return
}

// ListFriend implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFriend(ctx context.Context, req *relation.ListFriendReq) (resp *relation.ListFriendRes, err error) {
	resp = &relation.ListFriendRes{}
	userId := req.UserId

	followList, err := s.svcCtx.RelationService.ListFollow(req.UserId)
	if err != nil {
		return nil, err
	}
	followerList, err := s.svcCtx.RelationService.ListFollower(req.UserId)
	if err != nil {
		return nil, err
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
		}); err != nil {
			log.Println("get user info failed:", err)
		} else if res.Status != user.Status_OK {
			log.Println("user not exist")
			resp.Status = relation.Status_ERROR
			resp.ErrMsg = "user not exist"
		} else {
			isFollow, err := s.svcCtx.RelationService.IsFollow(req.UserId, userId)
			if err != nil {
				return nil, err
			}
			followCount, err := s.svcCtx.RelationService.GetFollowCount(userId)
			if err != nil {
				return nil, err
			}
			followerCount, err := s.svcCtx.RelationService.GetFollowerCount(userId)
			if err != nil {
				return nil, err
			}
			userInfo := convert2RelationUserInfo(res.User)
			userInfo.IsFollow = isFollow
			userInfo.FollowCount = followCount
			userInfo.FollowerCount = followerCount
			resp.Users[i] = &relation.FriendUser{
				User: userInfo,
			}
			if hasMessages {
				resp.Users[i].Message = messages[i].LastMessage
				resp.Users[i].MsgType = messages[i].MessageType
			}
		}
	}

	return
}

func convert2RelationUserInfo(user *user.UserInfo) *relation.UserInfo {
	return &relation.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
