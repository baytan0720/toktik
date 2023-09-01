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
	"toktik/internal/gateway/pkg/jwtutil"
	"toktik/internal/video/kitex_gen/video"
	mock_publishservice "toktik/pkg/test/mock/publish"
)

func newMockPublishClient(t *testing.T) *mock_publishservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_publishservice.NewMockClient(ctl)
}

func TestUserAPI_Publish(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockPublishClient(t)
		api := &PublishApi{
			publishClient: mockUserClient,
		}

		mockUserClient.EXPECT().PublishVideo(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		j := jwtutil.NewJwtUtil()
		token, _ := j.GenerateToken(jwtutil.CreateClaims(10))
		reqBody := PublishReq{
			Title: "haha",
			Data:  make([]byte, 1024),
			Token: token,
		}
		data, _ := json.Marshal(reqBody)
		ctx.Request.SetBodyRaw(data)
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())
		payload := PublishResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("publish failed", func(t *testing.T) {
		mockUserClient := newMockPublishClient(t)
		api := &PublishApi{
			publishClient: mockUserClient,
		}

		mockUserClient.EXPECT().PublishVideo(gomock.Any(), gomock.Any()).Return(&video.PublishVideoRes{
			Status: video.Status_ERROR,
			ErrMsg: PublishFail,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		j := jwtutil.NewJwtUtil()
		token, _ := j.GenerateToken(jwtutil.CreateClaims(10))
		reqBody := PublishReq{
			Title: "haha",
			Data:  make([]byte, 1024),
			Token: token,
		}
		data, _ := json.Marshal(reqBody)
		ctx.Request.SetBodyRaw(data)
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response.StatusCode())
		payload := PublishResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, PublishFail, payload.StatusMsg)
	})

	t.Run("publish success", func(t *testing.T) {
		mockUserClient := newMockPublishClient(t)
		api := &PublishApi{
			publishClient: mockUserClient,
		}

		mockUserClient.EXPECT().PublishVideo(gomock.Any(), gomock.Any()).Return(&video.PublishVideoRes{
			Status: video.Status_OK,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		j := jwtutil.NewJwtUtil()
		token, _ := j.GenerateToken(jwtutil.CreateClaims(10))
		reqBody := PublishReq{
			Title: "haha",
			Data:  make([]byte, 1024),
			Token: token,
		}
		data, _ := json.Marshal(reqBody)
		ctx.Request.SetBodyRaw(data)
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := PublishResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
	})
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

		assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())
		payload := ListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("list failed", func(t *testing.T) {
		mockUserClient := newMockPublishClient(t)
		api := &PublishApi{
			publishClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListVideo(gomock.Any(), gomock.Any()).Return(&video.ListVideoRes{
			Status: video.Status_ERROR,
			ErrMsg: PublishFail,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response.StatusCode())
		payload := PublishResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, PublishFail, payload.StatusMsg)
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
		payload := ListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
	})
}
