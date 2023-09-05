package relation

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/discovery"

	"toktik/internal/gateway/middleware"
	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/relation/kitex_gen/relation"
	"toktik/internal/relation/kitex_gen/relation/relationservice"
)

type RelationApi struct {
	relationClient relationservice.Client
}

func NewRelationApi(r discovery.Resolver) *RelationApi {
	return &RelationApi{
		relationClient: relationservice.MustNewClient("relation", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}

func (api *RelationApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/relation/action/",
			Handler: api.Action,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/relation/follow/list/",
			Handler: api.FollowList,
			Hooks:   []app.HandlerFunc{middleware.SoftAuthCheck},
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/relation/follower/list/",
			Handler: api.FollowerList,
			Hooks:   []app.HandlerFunc{middleware.SoftAuthCheck},
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/relation/friend/list/",
			Handler: api.FriendList,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
	}
}

type ActionResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func (api *RelationApi) Action(c context.Context, ctx *app.RequestContext) {
	actionType, err := strconv.Atoi(ctx.Query("action_type"))
	if err != nil {
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}
	toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}
	switch actionType {
	case 1:
		resp, err := api.relationClient.Follow(c, &relation.FollowReq{
			UserId:   ctx.GetInt64(middleware.CTX_USER_ID),
			ToUserId: toUserId,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		} else if resp.Status != 0 {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  resp.ErrMsg,
			})
			return
		}
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusOK,
		})
	case 2:
		resp, err := api.relationClient.Unfollow(c, &relation.UnfollowReq{
			UserId:   ctx.GetInt64(middleware.CTX_USER_ID),
			ToUserId: toUserId,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		} else if resp.Status != 0 {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  resp.ErrMsg,
			})
			return
		}
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusOK,
		})
	default:
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  "action_type error",
		})
	}

}

type FollowListResp struct {
	StatusCode int                  `json:"status_code"`
	StatusMsg  string               `json:"status_msg"`
	UserList   []*relation.UserInfo `json:"user_list"`
}

func (api *RelationApi) FollowList(c context.Context, ctx *app.RequestContext) {
	userId := ctx.GetInt64("user_id")
	resp, err := api.relationClient.ListFollow(c, &relation.ListFollowReq{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &FollowListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &FollowListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, &FollowListResp{
		StatusCode: apiutil.StatusOK,
		UserList:   resp.Users,
	})
}

type FollowerListResp struct {
	StatusCode int                  `json:"status_code"`
	StatusMsg  string               `json:"status_msg"`
	UserList   []*relation.UserInfo `json:"user_list"`
}

func (api *RelationApi) FollowerList(c context.Context, ctx *app.RequestContext) {
	userId := ctx.GetInt64("user_id")
	resp, err := api.relationClient.ListFollower(c, &relation.ListFollowerReq{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &FollowerListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &FollowerListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, &FollowerListResp{
		StatusCode: apiutil.StatusOK,
		UserList:   resp.Users,
	})
}

type FriendListResp struct {
	StatusCode int                    `json:"status_code"`
	StatusMsg  string                 `json:"status_msg"`
	UserList   []*relation.FriendUser `json:"user_list"`
}

func (api *RelationApi) FriendList(c context.Context, ctx *app.RequestContext) {
	userId := ctx.GetInt64("user_id")
	resp, err := api.relationClient.ListFriend(c, &relation.ListFriendReq{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &FriendListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &FriendListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, &FriendListResp{
		StatusCode: apiutil.StatusOK,
		UserList:   resp.Users,
	})
}
