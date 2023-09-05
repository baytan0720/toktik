package favorite

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"toktik/internal/favorite/kitex_gen/favorite"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"toktik/internal/gateway/pkg/apiutil"
	mock_favoriteservice "toktik/pkg/test/mock/favorite"
)

func newMockFavoriteClient(t *testing.T) *mock_favoriteservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_favoriteservice.NewMockClient(ctl)
}

func TestUserAPI_Favorite(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockFavoriteClient := newMockFavoriteClient(t)
		api := &FavoriteApi{
			favoriteClient: mockFavoriteClient,
		}

		mockFavoriteClient.EXPECT().Favorite(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10&action_type=1")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("favorite failed", func(t *testing.T) {
		mockUserClient := newMockFavoriteClient(t)
		api := &FavoriteApi{
			favoriteClient: mockUserClient,
		}

		mockUserClient.EXPECT().Favorite(gomock.Any(), gomock.Any()).Return(&favorite.FavoriteRes{
			Status: favorite.Status_ERROR,
			ErrMsg: FavoriteFail,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10&action_type=1")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, FavoriteFail, payload.StatusMsg)
	})

	t.Run("favorite success", func(t *testing.T) {
		mockUserClient := newMockFavoriteClient(t)
		api := &FavoriteApi{
			favoriteClient: mockUserClient,
		}

		mockUserClient.EXPECT().Favorite(gomock.Any(), gomock.Any()).Return(&favorite.FavoriteRes{
			Status: favorite.Status_OK,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10&action_type=1")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
	})
}

func TestUserAPI_List(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockFavoriteClient := newMockFavoriteClient(t)
		api := &FavoriteApi{
			favoriteClient: mockFavoriteClient,
		}

		mockFavoriteClient.EXPECT().ListFavorite(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("list failed", func(t *testing.T) {
		mockUserClient := newMockFavoriteClient(t)
		api := &FavoriteApi{
			favoriteClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListFavorite(gomock.Any(), gomock.Any()).Return(&favorite.ListFavoriteRes{
			Status: favorite.Status_ERROR,
			ErrMsg: FavoriteFail,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("user_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, FavoriteFail, payload.StatusMsg)
	})

	t.Run("List success", func(t *testing.T) {
		mockUserClient := newMockFavoriteClient(t)
		api := &FavoriteApi{
			favoriteClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListFavorite(gomock.Any(), gomock.Any()).Return(&favorite.ListFavoriteRes{
			Status: favorite.Status_OK,
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
