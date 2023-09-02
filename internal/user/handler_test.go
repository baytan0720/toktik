package main

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/internal/user/kitex_gen/user"
	"toktik/internal/user/pkg/ctx"
	usersvc "toktik/internal/user/pkg/user"
	mock_favoriteservice "toktik/pkg/test/mock/favorite"
	mock_relationservice "toktik/pkg/test/mock/relation"
	mock_videoservice "toktik/pkg/test/mock/video"
	"toktik/pkg/test/testutil"
)

func newMockUserService(t *testing.T) *usersvc.UserService {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&usersvc.User{}))
	return usersvc.NewUserService(func() *gorm.DB {
		return db
	})
}

func newMockRelationClient(t *testing.T) *mock_relationservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_relationservice.NewMockClient(ctl)
}

func newMockFavoriteClient(t *testing.T) *mock_favoriteservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_favoriteservice.NewMockClient(ctl)
}

func newMockVideoClient(t *testing.T) *mock_videoservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_videoservice.NewMockClient(ctl)
}

func TestUserServiceImpl_Register(t *testing.T) {
	t.Run("register failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			UserService: usersvc.NewUserService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewUserServiceImpl(svcCtx)
		_, err := svc.Register(context.Background(), &user.RegisterReq{
			Username: "123456",
			Password: "123456",
		})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			UserService: usersvc.NewUserService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewUserServiceImpl(svcCtx)
		resp, err := svc.Register(context.Background(), &user.RegisterReq{
			Username: "123456",
			Password: "123456",
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.UserId, "123456")
	})
}

func TestUserServiceImpl_Login(t *testing.T) {
	t.Run("login failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			UserService: usersvc.NewUserService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewUserServiceImpl(svcCtx)
		_, err := svc.Login(context.Background(), &user.LoginReq{
			Username: "123456",
			Password: "123456",
		})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			UserService: usersvc.NewUserService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewUserServiceImpl(svcCtx)
		resp, err := svc.Login(context.Background(), &user.LoginReq{
			Username: "123456",
			Password: "123456",
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.UserId, "123456")
	})
}

func TestUserServiceImpl_GetUserInfo(t *testing.T) {
	t.Run("login failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			UserService: usersvc.NewUserService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewUserServiceImpl(svcCtx)
		_, err := svc.GetUserInfo(context.Background(), &user.GetUserInfoReq{
			UserId: 1,
		})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRelationClient := newMockRelationClient(t)
		mockRelationClient.EXPECT().GetFollowInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		mockFavoriteClient := newMockFavoriteClient(t)
		mockFavoriteClient.EXPECT().GetFavoriteInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			UserService:    newMockUserService(t),
			RelationClient: mockRelationClient,
			FavoriteClient: mockFavoriteClient,
			VideoClient:    mockVideoClient,
		}
		svc := NewUserServiceImpl(svcCtx)
		resp, err := svc.GetUserInfo(context.Background(), &user.GetUserInfoReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.User, user.UserInfo{})
	})
}
