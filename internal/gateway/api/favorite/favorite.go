package favorite

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"

	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/favorite/kitex_gen/favorite/favoriteservice"
	"toktik/internal/gateway/middleware"
	"toktik/internal/gateway/pkg/apiutil"
)

var (
	FavoriteFail = "fail to favorite"
)

type FavoriteApi struct {
	favoriteClient favoriteservice.Client
}

func NewFavoriteApi(r discovery.Resolver) *FavoriteApi {
	return &FavoriteApi{
		favoriteClient: favoriteservice.MustNewClient("favorite", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}

func (api *FavoriteApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/favorite/action",
			Handler: api.Action,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/favorite/list",
			Handler: api.List,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
	}
}

type FavoriteResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func (api *FavoriteApi) Action(c context.Context, ctx *app.RequestContext) {
	actionType, err := strconv.Atoi(ctx.Query("action_type"))
	videoId, err1 := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err != nil || err1 != nil {
		ctx.JSON(http.StatusBadRequest, &FavoriteResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}

	switch actionType {
	case 1:
		resp, err := api.favoriteClient.Favorite(c, &favorite.FavoriteReq{
			VideoId: videoId,
			UserId:  ctx.GetInt64(middleware.CTX_USER_ID),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &FavoriteResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		} else if resp.Status != 0 {
			ctx.JSON(http.StatusBadRequest, &FavoriteResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  resp.ErrMsg,
			})
			return
		}
		break
	case 2:
		resp, err := api.favoriteClient.UnFavorite(c, &favorite.UnFavoriteReq{
			VideoId: videoId,
			UserId:  ctx.GetInt64(middleware.CTX_USER_ID),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &FavoriteResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		} else if resp.Status != 0 {
			ctx.JSON(http.StatusBadRequest, &FavoriteResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  resp.ErrMsg,
			})
			return
		}
		break
	default:
		ctx.JSON(http.StatusBadRequest, &FavoriteResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &FavoriteResp{
		StatusCode: apiutil.StatusOK,
	})
}

type ListResp struct {
	StatusCode int
	StatusMsg  string
	VideoList  []*favorite.VideoInfo
}

func (api *FavoriteApi) List(c context.Context, ctx *app.RequestContext) {
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &ListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}

	resp, err := api.favoriteClient.ListFavorite(c, &favorite.ListFavoriteReq{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &ListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusBadRequest, &ListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, &ListResp{
		StatusCode: apiutil.StatusOK,
		VideoList:  resp.VideoList,
	})
}
