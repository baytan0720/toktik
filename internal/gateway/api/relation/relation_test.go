package relation

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
	"toktik/internal/relation/kitex_gen/relation"
	mock_relationservice "toktik/pkg/test/mock/relation"
)

func newMockUserClient(t *testing.T) *mock_relationservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_relationservice.NewMockClient(ctl)
}

func TestUserAPI_Action(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}

		mockUserClient.EXPECT().Follow(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10&action_type=1")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("follow failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}

		mockUserClient.EXPECT().Follow(gomock.Any(), gomock.Any()).Return(&relation.FollowRes{
			Status: relation.Status_ERROR,
			ErrMsg: "fail to follow",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10&action_type=1")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "fail to follow", payload.StatusMsg)
	})

	t.Run("follow success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}

		mockUserClient.EXPECT().Follow(gomock.Any(), gomock.Any()).Return(&relation.FollowRes{
			Status: relation.Status_OK,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("to_user_id=10&action_type=1")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
	})
}

func TestUserAPI_FollowList(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListFollow(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.FollowList(context.Background(), ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())
		payload := FollowListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("follow list failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListFollow(gomock.Any(), gomock.Any()).Return(&relation.ListFollowRes{
			Status: relation.Status_ERROR,
			ErrMsg: "fail to follow",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.FollowList(context.Background(), ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response.StatusCode())
		payload := FollowListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "fail to follow", payload.StatusMsg)
	})

	t.Run("follow list success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}
		infos := make([]*relation.UserInfo, 10)
		mockUserClient.EXPECT().ListFollow(gomock.Any(), gomock.Any()).Return(&relation.ListFollowRes{
			Status: relation.Status_OK,
			Users:  infos,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.FollowList(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := FollowListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, infos, payload.UserList)
	})
}

func TestUserAPI_FollowerList(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListFollower(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.FollowerList(context.Background(), ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())
		payload := FollowerListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("follower list failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListFollower(gomock.Any(), gomock.Any()).Return(&relation.ListFollowerRes{
			Status: relation.Status_ERROR,
			ErrMsg: "fail to follow",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.FollowerList(context.Background(), ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response.StatusCode())
		payload := FollowerListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "fail to follow", payload.StatusMsg)
	})

	t.Run("follower list success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		api := &RelationApi{
			relationClient: mockUserClient,
		}
		infos := make([]*relation.UserInfo, 10)
		mockUserClient.EXPECT().ListFollower(gomock.Any(), gomock.Any()).Return(&relation.ListFollowerRes{
			Status: relation.Status_OK,
			Users:  infos,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.FollowerList(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := FollowerListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, infos, payload.UserList)
	})
}
