package main

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/internal/video/kitex_gen/video"
	"toktik/internal/video/pkg/ctx"
	videosvc "toktik/internal/video/pkg/video"
	mock_userservice "toktik/pkg/test/mock/user"
	"toktik/pkg/test/testutil"
)

func newMockVideoService(t *testing.T) *videosvc.VideoService {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&videosvc.Video{}))
	return videosvc.NewVideoService(func() *gorm.DB {
		return db
	})
}

func newMockUserClient(t *testing.T) *mock_userservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_userservice.NewMockClient(ctl)
}

func TestVideoServiceImpl_GetVideo(t *testing.T) {
	t.Run("get relation info failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			VideoService: videosvc.NewVideoService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewVideoServiceImpl(svcCtx)
		_, err := svc.GetVideo(context.Background(), &video.GetVideoReq{
			VideoId: 1,
		})
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(&video.GetVideoRes{Video: &video.Video{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			VideoService: newMockVideoService(t),
			VideoClient:  mockVideoClient,
		}
		svc := NewVideoServiceImpl(svcCtx)
		resp, err := svc.GetVideo(context.Background(), &video.GetVideoReq{
			VideoId: 1,
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
func TestVideoServiceImpl_ListVideo(t *testing.T) {
	t.Run("list video failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			VideoService: videosvc.NewVideoService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewVideoServiceImpl(svcCtx)
		_, err := svc.ListVideo(context.Background(), &video.ListVideoReq{
			UserId: 1,
		})
		assert.Error(t, err)
	})
	t.Run("get video info failed", func(t *testing.T) {
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			VideoService: newMockVideoService(t),
			VideoClient:  mockVideoClient,
		}
		svc := NewVideoServiceImpl(svcCtx)
		resp, err := svc.ListVideo(context.Background(), &video.ListVideoReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("success", func(t *testing.T) {
		mockVideoClient := newMockVideoClient(t)
		mockVideoClient.EXPECT().GetVideo(gomock.Any(), gomock.Any()).Return(&video.GetVideoRes{Video: &video.Video{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			VideoService: newMockVideoService(t),
			VideoClient:  mockVideoClient,
		}
		svc := NewVideoServiceImpl(svcCtx)
		resp, err := svc.ListVideo(context.Background(), &video.ListVideoReq{
			UserId: 1,
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestVideoServiceImpl_PublishVideo(t *testing.T) {
	t.Run("publish video failed", func(t *testing.T) {
		// 创建一个模拟的视频上传客户端
		mockVideoClient := newMockVideoClient(t)
		// 设置模拟视频上传客户端的PublishVideo方法期望，使其返回错误
		mockVideoClient.EXPECT().PublishVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		// 创建一个模拟的服务上下文
		svcCtx := &ctx.ServiceContext{
			VideoClient: mockVideoClient,
		}
		// 创建一个模拟的视频服务实现
		svc := &VideoServiceImpl{
			svcCtx: svcCtx,
		}
		// 调用PublishVideo方法，并传入一个模拟的PublishVideo请求
		_, err := svc.PublishVideo(context.Background(), &video.PublishVideoReq{}, &video.Video{})
		// 检查是否返回了错误
		assert.Error(t, err)
	})
}
