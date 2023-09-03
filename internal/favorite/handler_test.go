package main

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/favorite/pkg/ctx"
	favoritesvc "toktik/internal/favorite/pkg/favorite"
	"toktik/internal/video/kitex_gen/video"
	"toktik/pkg/db/model"
	mock_videoservice "toktik/pkg/test/mock/video"
	"toktik/pkg/test/testutil"
)

func newMockFavoriteService(t *testing.T) *favoritesvc.FavoriteService {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&favoritesvc.Favorite{}))
	return favoritesvc.NewFavoriteService(func() *gorm.DB {
		return db
	})
}

func newMockVideoClient(t *testing.T) *mock_videoservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_videoservice.NewMockClient(ctl)
}

func TestFavoriteServiceImpl_Favorite(t *testing.T) {
	t.Run("favorite failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(&video.GetVideoRes{VideoList: []*video.Video{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: favoritesvc.NewFavoriteService(func() *gorm.DB {
				return db
			}),
			VideoClient: mockVideoClient,
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.Favorite(context.Background(), &favorite.FavoriteReq{
			UserId:  1,
			VideoId: 1,
		})
		assert.Error(t, err)
	})

	t.Run("get video info failed", func(t *testing.T) {
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: newMockFavoriteService(t),
			VideoClient:     mockVideoClient,
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.Favorite(context.Background(), &favorite.FavoriteReq{
			UserId:  1,
			VideoId: 2,
		})
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(&video.GetVideoRes{VideoList: []*video.Video{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: newMockFavoriteService(t),
			VideoClient:     mockVideoClient,
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.Favorite(context.Background(), &favorite.FavoriteReq{
			UserId:  1,
			VideoId: 2,
		})
		assert.NoError(t, err)
	})
}
func TestUnFavoriteServiceImpl_UnFavorite(t *testing.T) {
	t.Run("unfavorite failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: favoritesvc.NewFavoriteService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.UnFavorite(context.Background(), &favorite.UnFavoriteReq{
			UserId:  1,
			VideoId: 1,
		})
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		db := testutil.NewMockDB()
		require.NoError(t, db.AutoMigrate(&favoritesvc.Favorite{}))

		db.Create(&model.Favorite{
			UserId:     1,
			VideoId:    1,
			IsFavorite: true,
		})

		svcCtx := &ctx.ServiceContext{
			FavoriteService: favoritesvc.NewFavoriteService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.UnFavorite(context.Background(), &favorite.UnFavoriteReq{
			UserId:  1,
			VideoId: 1,
		})
		assert.NoError(t, err)
	})
}

func TestFavoriteServiceImpl_ListFavorite(t *testing.T) {
	t.Run("list favorite failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: favoritesvc.NewFavoriteService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.ListFavorite(context.Background(), &favorite.ListFavoriteReq{
			UserId: 1,
		})
		assert.Error(t, err)
	})

	t.Run("get video info failed", func(t *testing.T) {
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: newMockFavoriteService(t),
			VideoClient:     mockVideoClient,
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		resp, err := svc.ListFavorite(context.Background(), &favorite.ListFavoriteReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.VideoList))
		assert.NotNil(t, resp)
	})

	t.Run("success", func(t *testing.T) {
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(&video.GetVideoRes{VideoList: []*video.Video{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: newMockFavoriteService(t),
			VideoClient:     mockVideoClient,
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.ListFavorite(context.Background(), &favorite.ListFavoriteReq{
			UserId: 1,
		})
		assert.NoError(t, err)
	})
}

func TestFavoriteServiceImpl_GetFavoriteInfo(t *testing.T) {
	t.Run("get favorite info failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			FavoriteService: favoritesvc.NewFavoriteService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		_, err := svc.GetFavoriteInfo(context.Background(), &favorite.GetFavoriteInfoReq{
			UserId:      1,
			VideoIdList: []int64{1},
		})
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		svcCtx := &ctx.ServiceContext{
			FavoriteService: newMockFavoriteService(t),
		}
		svc := NewFavoriteServiceImpl(svcCtx)
		resp, err := svc.GetFavoriteInfo(context.Background(), &favorite.GetFavoriteInfoReq{
			UserId:      1,
			VideoIdList: []int64{2},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.FavoriteInfoList))
		assert.Equal(t, int64(0), resp.FavoriteInfoList[0].Count)
	})
}
