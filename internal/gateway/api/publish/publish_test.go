package publish

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/video/kitex_gen/video"
	mock_videoservice "toktik/pkg/test/mock/video"
)

func newMockPublishClient(t *testing.T) *mock_videoservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_videoservice.NewMockClient(ctl)
}

func TestUserAPI_List(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockPublishClient(t)
		api := &PublishApi{
			publishClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListVideo(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ListRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, apiutil.ErrInternalError.Error(), payload.StatusMsg)
	})

	t.Run("list success", func(t *testing.T) {
		mockUserClient := newMockPublishClient(t)
		api := &PublishApi{
			publishClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListVideo(gomock.Any(), gomock.Any()).Return(&video.ListVideoRes{
			Status: video.Status_OK,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ListRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
	})
}
