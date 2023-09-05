package feed

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/video/kitex_gen/video"
	mock_videoservice "toktik/pkg/test/mock/video"
)

func newMockFeedClient(t *testing.T) *mock_videoservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_videoservice.NewMockClient(ctl)
}

func TestUserAPI_Feed(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockFeedClient := newMockFeedClient(t)
		api := &FeedApi{
			videoClient: mockFeedClient,
		}

		mockFeedClient.EXPECT().Feed(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("latest_time=" + strconv.FormatInt(time.Now().UnixMilli(), 10))
		api.Feed(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := FeedResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("get feed failed", func(t *testing.T) {
		mockUserClient := newMockFeedClient(t)
		api := &FeedApi{
			videoClient: mockUserClient,
		}

		mockUserClient.EXPECT().Feed(gomock.Any(), gomock.Any()).Return(&video.FeedRes{
			Status: video.Status_ERROR,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("latest_time='string'")
		api.Feed(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := FeedResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
	})

	t.Run("get feed success", func(t *testing.T) {
		mockFeedClient := newMockFeedClient(t)
		api := &FeedApi{
			videoClient: mockFeedClient,
		}

		infos := make([]*video.VideoInfo, 5)

		mockFeedClient.EXPECT().Feed(gomock.Any(), gomock.Any()).Return(&video.FeedRes{
			Status:    video.Status_OK,
			NextTime:  time.Now().UnixMilli(),
			VideoList: infos,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("latest_time=" + strconv.FormatInt(time.Now().UnixMilli(), 10))
		api.Feed(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := FeedResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, infos, payload.VideoList)
		assert.NotEmpty(t, payload.NextTime)
	})
}
