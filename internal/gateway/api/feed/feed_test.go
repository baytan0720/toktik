package feed

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"toktik/internal/feed/kitex_gen/feed"
	"toktik/internal/gateway/pkg/apiutil"
	mock_feedservice "toktik/pkg/test/mock/feed"
)

func newMockFeedClient(t *testing.T) *mock_feedservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_feedservice.NewMockClient(ctl)
}

func TestUserAPI_Feed(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockFeedClient := newMockFeedClient(t)
		api := &FeedApi{
			feedClient: mockFeedClient,
		}

		mockFeedClient.EXPECT().Feed(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("latest_time='string'")
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
			feedClient: mockUserClient,
		}

		mockUserClient.EXPECT().Feed(gomock.Any(), gomock.Any()).Return(&feed.FeedRes{
			Status: feed.Status_ERROR,
			ErrMsg: FeedFail,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("latest_time='string'")
		api.Feed(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := FeedResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, FeedFail, payload.StatusMsg)
	})

	t.Run("get feed success", func(t *testing.T) {
		mockFeedClient := newMockFeedClient(t)
		api := &FeedApi{
			feedClient: mockFeedClient,
		}

		infos := make([]*feed.VideoInfo, 5)

		mockFeedClient.EXPECT().Feed(gomock.Any(), gomock.Any()).Return(&feed.FeedRes{
			Status:    feed.Status_OK,
			NextTime:  "nexttime",
			VideoList: infos,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("username=test_user&password=123456")
		api.Feed(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := FeedResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, infos, payload.VideoList)
		assert.NotEmpty(t, payload.NextTime)
	})
}
