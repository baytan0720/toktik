package user

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
	"toktik/internal/user/kitex_gen/user"
	mock_userservice "toktik/pkg/test/mock/user"
)

func newMockUserClient(t *testing.T) *mock_userservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_userservice.NewMockClient(ctl)
}

func TestUserAPI_Register(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		ctx := app.NewContext(16)
		api.Register(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := RegisterRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, apiutil.ErrInvalidParams.Error(), payload.StatusMsg)
	})

	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("username=test_user&password=123456")
		api.Register(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := LoginRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, apiutil.ErrInternalError.Error(), payload.StatusMsg)
	})

	t.Run("Register failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&user.RegisterRes{
			Status: user.Status_ERROR,
			ErrMsg: "用户名被占用",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("username=test_user&password=123456")
		api.Register(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := LoginRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "用户名被占用", payload.StatusMsg)
	})

	t.Run("Register success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&user.RegisterRes{
			Status: user.Status_OK,
			UserId: 10,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("username=test_user&password=123456")
		api.Register(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := RegisterRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, int64(10), payload.UserId)
		assert.NotEmpty(t, payload.Token)
	})
}

func TestUserAPI_Login(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		ctx := app.NewContext(16)
		api.Login(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := LoginRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, apiutil.ErrInvalidParams.Error(), payload.StatusMsg)
	})

	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("username=test_user&password=123456")
		api.Login(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := LoginRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, apiutil.ErrInternalError.Error(), payload.StatusMsg)
	})

	t.Run("login failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&user.LoginRes{
			Status: user.Status_ERROR,
			ErrMsg: "用户名或密码错误",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("username=test_user&password=123456")
		api.Login(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := LoginRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "用户名或密码错误", payload.StatusMsg)
	})

	t.Run("login success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&user.LoginRes{
			Status: user.Status_OK,
			UserId: 10,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("username=test_user&password=123456")
		api.Login(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := LoginRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, int64(10), payload.UserId)
		assert.NotEmpty(t, payload.Token)
	})
}

func TestUserAPI_GetUserInfo(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		ctx := app.NewContext(16)
		api.UserInfo(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := UserInfoRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, apiutil.ErrInvalidParams.Error(), payload.StatusMsg)
	})

	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.UserInfo(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := UserInfoRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, apiutil.ErrInternalError.Error(), payload.StatusMsg)
	})

	t.Run("get info failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{
			Status: user.Status_ERROR,
			ErrMsg: "get user info failed",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.UserInfo(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := UserInfoRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "get user info failed", payload.StatusMsg)
	})

	t.Run("get info success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &UserAPI{
			userClient: mockUserClient,
		}

		userinfo := &user.UserInfo{
			Id:   10,
			Name: "nice",
		}

		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{
			Status: user.Status_OK,
			User:   userinfo,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10&token='string'")
		api.UserInfo(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := UserInfoRes{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, userinfo, payload.Info)
	})
}
