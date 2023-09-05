package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/relation/kitex_gen/relation"
	"toktik/internal/user/kitex_gen/user"
	"toktik/internal/user/pkg/ctx"
	usersvc "toktik/internal/user/pkg/user"
	"toktik/internal/video/kitex_gen/video"
	mock_favoriteservice "toktik/pkg/test/mock/favorite"
	mock_relationservice "toktik/pkg/test/mock/relation"
	mock_videoservice "toktik/pkg/test/mock/video"
	"toktik/pkg/test/testutil"
)

func newMockRelationClient(t *testing.T) *mock_relationservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_relationservice.NewMockClient(ctl)
}

func newMockVideoClient(t *testing.T) *mock_videoservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_videoservice.NewMockClient(ctl)
}

func newMockFavoriteClient(t *testing.T) *mock_favoriteservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_favoriteservice.NewMockClient(ctl)
}

func TestUserServiceImpl_Register(t *testing.T) {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&usersvc.User{}))
	userSvc := usersvc.NewUserService(func() *gorm.DB {
		return db
	})
	svcCtx := &ctx.ServiceContext{
		UserService: userSvc,
	}
	svc := NewUserServiceImpl(svcCtx)

	t.Run("success", func(t *testing.T) {
		resp, err := svc.Register(context.Background(), &user.RegisterReq{
			Username: "test",
			Password: "123456",
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_OK, resp.Status)
	})

	t.Run("duplicate username", func(t *testing.T) {
		resp, err := svc.Register(context.Background(), &user.RegisterReq{
			Username: "test",
			Password: "123456",
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_ERROR, resp.Status)
	})

	t.Run("password too short", func(t *testing.T) {
		resp, err := svc.Register(context.Background(), &user.RegisterReq{
			Username: "test2",
			Password: "123",
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_ERROR, resp.Status)
	})
}

func TestUserServiceImpl_Login(t *testing.T) {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&usersvc.User{}))
	userSvc := usersvc.NewUserService(func() *gorm.DB {
		return db
	})
	svcCtx := &ctx.ServiceContext{
		UserService: userSvc,
	}
	svc := NewUserServiceImpl(svcCtx)

	t.Run("user not exist", func(t *testing.T) {
		resp, err := svc.Login(context.Background(), &user.LoginReq{
			Username: "test",
			Password: "123456",
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_ERROR, resp.Status)
	})

	err := db.Create(&usersvc.User{
		Username: "test",
		Password: fmt.Sprintf("%x", md5.Sum([]byte("123456"))),
	}).Error
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		resp, err := svc.Login(context.Background(), &user.LoginReq{
			Username: "test",
			Password: "123456",
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_OK, resp.Status)
	})

	t.Run("wrong password", func(t *testing.T) {
		resp, err := svc.Login(context.Background(), &user.LoginReq{
			Username: "test",
			Password: "1234567",
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_ERROR, resp.Status)
	})
}

func TestUserServiceImpl_GetUserInfo(t *testing.T) {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&usersvc.User{}))
	userSvc := usersvc.NewUserService(func() *gorm.DB {
		return db
	})
	mockFavoriteClient := newMockFavoriteClient(t)
	mockVideoClient := newMockVideoClient(t)
	mockRelationClient := newMockRelationClient(t)

	svcCtx := &ctx.ServiceContext{
		UserService:    userSvc,
		FavoriteClient: mockFavoriteClient,
		VideoClient:    mockVideoClient,
		RelationClient: mockRelationClient,
	}
	svc := NewUserServiceImpl(svcCtx)

	t.Run("user not exist", func(t *testing.T) {
		resp, err := svc.GetUserInfo(context.Background(), &user.GetUserInfoReq{
			ToUserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_ERROR, resp.Status)
	})

	err := db.Create(&usersvc.User{
		Username: "test",
		Password: fmt.Sprintf("%x", md5.Sum([]byte("123456"))),
	}).Error
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockFavoriteClient.EXPECT().GetUserFavoriteInfo(gomock.Any(), gomock.Any()).Return(&favorite.GetUserFavoriteInfoRes{
			FavoriteInfoList: []*favorite.UserFavoriteInfo{
				{
					UserId:         1,
					FavoriteCount:  100,
					TotalFavorited: 200,
				},
			},
		}, nil)

		mockRelationClient.EXPECT().GetFollowInfo(gomock.Any(), gomock.Any()).Return(&relation.GetFollowInfoRes{
			FollowInfoList: []*relation.FollowInfo{
				{
					UserId:        1,
					FollowCount:   100,
					FollowerCount: 200,
					IsFollow:      true,
				},
			},
		}, nil)

		mockVideoClient.EXPECT().GetWorkCount(gomock.Any(), gomock.Any()).Return(&video.GetWorkCountRes{
			WorkCountList: []*video.WorkCountInfo{
				{
					UserId:    1,
					WorkCount: 100,
				},
			},
		}, nil)

		resp, err := svc.GetUserInfo(context.Background(), &user.GetUserInfoReq{
			ToUserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, user.Status_OK, resp.Status)
		expected := user.UserInfo{
			Id:             1,
			Name:           "test",
			FollowCount:    100,
			FollowerCount:  200,
			IsFollow:       true,
			WorkCount:      100,
			FavoriteCount:  100,
			TotalFavorited: 200,
		}
		assert.Equal(t, expected.Id, resp.User.Id)
		assert.Equal(t, expected.Name, resp.User.Name)
		assert.Equal(t, expected.FollowCount, resp.User.FollowCount)
		assert.Equal(t, expected.FollowerCount, resp.User.FollowerCount)
		assert.Equal(t, expected.IsFollow, resp.User.IsFollow)
		assert.Equal(t, expected.WorkCount, resp.User.WorkCount)
		assert.Equal(t, expected.FavoriteCount, resp.User.FavoriteCount)
		assert.Equal(t, expected.TotalFavorited, resp.User.TotalFavorited)
	})
}
