package main

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/internal/message/kitex_gen/message"
	"toktik/internal/relation/kitex_gen/relation"
	"toktik/internal/relation/pkg/ctx"
	relationsvc "toktik/internal/relation/pkg/relation"
	"toktik/internal/user/kitex_gen/user"
	"toktik/pkg/db/model"
	mock_messageservice "toktik/pkg/test/mock/message"
	mock_userservice "toktik/pkg/test/mock/user"
	"toktik/pkg/test/testutil"
)

func newMockRelationService(t *testing.T) *relationsvc.RelationService {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&relationsvc.Relation{}))
	return relationsvc.NewRelationService(func() *gorm.DB {
		return db
	})
}

func newMockUserClient(t *testing.T) *mock_userservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_userservice.NewMockClient(ctl)
}

func newMockMessageClient(t *testing.T) *mock_messageservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_messageservice.NewMockClient(ctl)
}

func TestRelationServiceImpl_GetFollowInfo(t *testing.T) {
	t.Run("get relation info failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.GetFollowInfo(context.Background(), &relation.GetFollowInfoReq{
			UserId:       1,
			ToUserIdList: []int64{1},
		})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: newMockRelationService(t),
			UserClient:      mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		// ToUser
		resp, err := svc.GetFollowInfo(context.Background(), &relation.GetFollowInfoReq{
			UserId:       1,
			ToUserIdList: []int64{2},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.FollowInfoList))
		assert.Equal(t, int64(0), resp.FollowInfoList[0].FollowerCount)
		assert.Equal(t, int64(0), resp.FollowInfoList[0].FollowCount)
		assert.False(t, resp.FollowInfoList[0].IsFollow)
	})
}

func TestRelationServiceImpl_Follow(t *testing.T) {
	t.Run("follow failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
			UserClient: mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.Follow(context.Background(), &relation.FollowReq{
			UserId:   1,
			ToUserId: 2,
		})
		assert.Error(t, err)
	})

	t.Run("get user info failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: newMockRelationService(t),
			UserClient:      mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.Follow(context.Background(), &relation.FollowReq{
			UserId:   1,
			ToUserId: 2,
		})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: newMockRelationService(t),
			UserClient:      mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.Follow(context.Background(), &relation.FollowReq{
			UserId:   1,
			ToUserId: 2,
		})
		assert.NoError(t, err)
	})
}

func TestRelationServiceImpl_Unfollow(t *testing.T) {
	t.Run("unfollow failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.Unfollow(context.Background(), &relation.UnfollowReq{
			UserId:   1,
			ToUserId: 2,
		})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		db := testutil.NewMockDB()
		require.NoError(t, db.AutoMigrate(&relationsvc.Relation{}))

		db.Create(&model.Relation{
			UserId:   1,
			ToUserId: 2,
			IsFollow: false,
		})

		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB{
				return db
			}),
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.Unfollow(context.Background(), &relation.UnfollowReq{
			UserId:   1,
			ToUserId: 2,
		})
		assert.NoError(t, err)
	})
}

func TestRelationServiceImpl_ListFollow(t *testing.T) {
	t.Run("list follow failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.ListFollow(context.Background(), &relation.ListFollowReq{
			UserId: 1,
		})
		assert.Error(t, err)
	})

	t.Run("get user info failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: newMockRelationService(t),
			UserClient:      mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		resp, err := svc.ListFollow(context.Background(), &relation.ListFollowReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.Users))
		assert.NotNil(t, resp)
	})

	t.Run("success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: newMockRelationService(t),
			UserClient:      mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.ListFollow(context.Background(), &relation.ListFollowReq{
			UserId: 1,
		})
		assert.NoError(t, err)
	})
}

func TestRelationServiceImpl_ListFollower(t *testing.T) {
	t.Run("list follower failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.ListFollower(context.Background(), &relation.ListFollowerReq{
			UserId: 1,
		})
		assert.Error(t, err)
	})

	t.Run("get user info failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: newMockRelationService(t),
			UserClient:      mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		resp, err := svc.ListFollower(context.Background(), &relation.ListFollowerReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.Users))
		assert.NotNil(t, resp)
	})

	t.Run("success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			RelationService: newMockRelationService(t),
			UserClient:      mockUserClient,
		}
		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.ListFollower(context.Background(), &relation.ListFollowerReq{
			UserId: 1,
		})
		assert.NoError(t, err)
	})
}

func TestRelationServiceImpl_ListFriend(t *testing.T) {
	t.Run("list friend failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
		}

		svc := NewRelationServiceImpl(svcCtx)
		_, err := svc.ListFriend(context.Background(), &relation.ListFriendReq{
			UserId: 1,
		})
		assert.Error(t, err)
	})

	t.Run("get messages failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		require.NoError(t, db.AutoMigrate(&relationsvc.Relation{}))

		db.Create(&model.Relation{
			UserId:   1,
			ToUserId: 2,
			IsFollow: true,
		})
		db.Create(&model.Relation{
			UserId:   2,
			ToUserId: 1,
			IsFollow: true,
		})

		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		mockMessageClient := newMockMessageClient(t)
		mockMessageClient.EXPECT().GetLastMessage(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()

		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
			UserClient:    mockUserClient,
			MessageClient: mockMessageClient,
		}
		svc := NewRelationServiceImpl(svcCtx)

		resp, err := svc.ListFriend(context.Background(), &relation.ListFriendReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Users))
		assert.Equal(t, int64(0), resp.Users[0].MsgType)
		assert.Equal(t, "", resp.Users[0].Message)
	})

	t.Run("get user info failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		require.NoError(t, db.AutoMigrate(&relationsvc.Relation{}))

		db.Create(&model.Relation{
			UserId:   1,
			ToUserId: 2,
			IsFollow: true,
		})
		db.Create(&model.Relation{
			UserId:   2,
			ToUserId: 1,
			IsFollow: true,
		})

		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		mockMessageClient := newMockMessageClient(t)
		mockMessageClient.EXPECT().GetLastMessage(gomock.Any(), gomock.Any()).Return(&message.GetLastMessageRes{Messages: []*message.LastMessage{&message.LastMessage{
			LastMessage: "test message",
			MessageType: 0,
		}}}, nil).AnyTimes()

		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
			UserClient:    mockUserClient,
			MessageClient: mockMessageClient,
		}
		svc := NewRelationServiceImpl(svcCtx)

		resp, err := svc.ListFriend(context.Background(), &relation.ListFriendReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Users))
		assert.Nil(t, resp.Users[0])
	})

	t.Run("success", func(t *testing.T) {
		db := testutil.NewMockDB()
		require.NoError(t, db.AutoMigrate(&relationsvc.Relation{}))

		db.Create(&model.Relation{
			UserId:   1,
			ToUserId: 2,
			IsFollow: true,
		})
		db.Create(&model.Relation{
			UserId:   2,
			ToUserId: 1,
			IsFollow: true,
		})

		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		mockMessageClient := newMockMessageClient(t)
		mockMessageClient.EXPECT().GetLastMessage(gomock.Any(), gomock.Any()).Return(&message.GetLastMessageRes{Messages: []*message.LastMessage{&message.LastMessage{
			LastMessage: "test message",
			MessageType: 0,
		}}}, nil).AnyTimes()

		svcCtx := &ctx.ServiceContext{
			RelationService: relationsvc.NewRelationService(func() *gorm.DB {
				return db
			}),
			UserClient:    mockUserClient,
			MessageClient: mockMessageClient,
		}
		svc := NewRelationServiceImpl(svcCtx)

		resp, err := svc.ListFriend(context.Background(), &relation.ListFriendReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Users))
		assert.Equal(t, int64(0), resp.Users[0].MsgType)
		assert.Equal(t, "test message", resp.Users[0].Message)
	})

}
