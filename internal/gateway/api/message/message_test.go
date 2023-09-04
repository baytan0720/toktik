package message

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
	"toktik/internal/message/kitex_gen/message"
	mock_messageservice "toktik/pkg/test/mock/message"
)

func newMockMessageClient(t *testing.T) *mock_messageservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_messageservice.NewMockClient(ctl)
}

func TestUserAPI_Action(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockMessageClient(t)
		api := &MessageApi{
			messageClient: mockUserClient,
		}

		mockUserClient.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10&action_type=1&content='hello'")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("send message failed", func(t *testing.T) {
		mockUserClient := newMockMessageClient(t)
		api := &MessageApi{
			messageClient: mockUserClient,
		}

		mockUserClient.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(&message.SendMessageRes{
			Status: message.Status_ERROR,
			ErrMsg: "fail to send",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10&action_type=1&content='hello'")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "fail to send", payload.StatusMsg)
	})

	t.Run("send success", func(t *testing.T) {
		mockUserClient := newMockMessageClient(t)
		api := &MessageApi{
			messageClient: mockUserClient,
		}

		mockUserClient.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(&message.SendMessageRes{
			Status: message.Status_OK,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10&action_type=1&content='hello'")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
	})
}

func TestUserAPI_List(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockMessageClient(t)
		api := &MessageApi{
			messageClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListMessage(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10")
		api.Chat(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ChatResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("list failed", func(t *testing.T) {
		mockUserClient := newMockMessageClient(t)
		api := &MessageApi{
			messageClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListMessage(gomock.Any(), gomock.Any()).Return(&message.ListMessageRes{
			Status: message.Status_ERROR,
			ErrMsg: "fail to list",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10")
		api.Chat(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ChatResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "fail to list", payload.StatusMsg)
	})

	t.Run("list success", func(t *testing.T) {
		mockUserClient := newMockMessageClient(t)
		api := &MessageApi{
			messageClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListMessage(gomock.Any(), gomock.Any()).Return(&message.ListMessageRes{
			Status: message.Status_OK,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10")
		api.Chat(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ChatResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
	})
}