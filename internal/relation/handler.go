package main

import (
	"context"

	"toktik/internal/relation/kitex_gen/relation"
	"toktik/internal/relation/pkg/ctx"
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

// GetFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollow(ctx context.Context, req *relation.GetFollowInfoReq) (resp *relation.GetFollowInfoRes, err error) {
	// TODO: Your code here...
	return
}

// Follow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Follow(ctx context.Context, req *relation.FollowReq) (resp *relation.FollowRes, err error) {
	// TODO: Your code here...
	return
}

// Unfollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Unfollow(ctx context.Context, req *relation.UnfollowReq) (resp *relation.UnfollowRes, err error) {
	// TODO: Your code here...
	return
}

// ListFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollow(ctx context.Context, req *relation.ListFollowReq) (resp *relation.ListFollowRes, err error) {
	// TODO: Your code here...
	return
}

// ListFollower implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFollower(ctx context.Context, req *relation.ListFollowerReq) (resp *relation.ListFollowerRes, err error) {
	// TODO: Your code here...
	return
}

// ListFriend implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ListFriend(ctx context.Context, req *relation.ListFriendReq) (resp *relation.ListFriendRes, err error) {
	// TODO: Your code here...
	return
}
