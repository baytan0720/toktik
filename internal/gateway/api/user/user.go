package user

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"

	"toktik/internal/gateway/middleware"
	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/gateway/pkg/jwtutil"
	"toktik/internal/user/kitex_gen/user"
	"toktik/internal/user/kitex_gen/user/userservice"
)

type UserAPI struct {
	userClient userservice.Client
}

func NewUserAPI(r discovery.Resolver) *UserAPI {
	return &UserAPI{
		userClient: userservice.MustNewClient("user", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}

func (api *UserAPI) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/user/register/",
			Handler: api.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/douyin/user/login/",
			Handler: api.Login,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/user/",
			Handler: api.UserInfo,
			Hooks:   []app.HandlerFunc{middleware.SoftAuthCheck},
		},
	}
}

type LoginOrRegisterReq struct {
	Username string `query:"username,required"`
	Password string `query:"password,required"`
}

type RegisterRes struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

func (api *UserAPI) Register(c context.Context, ctx *app.RequestContext) {
	params := &LoginOrRegisterReq{}
	if err := ctx.Bind(params); apiutil.HandleError(ctx, err, apiutil.ErrInvalidParams) {
		return
	}

	resp, err := api.userClient.Register(c, &user.RegisterReq{
		Username: params.Username,
		Password: params.Password,
	})
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
		return
	}

	token, err := jwtutil.GenerateTokenWithUserId(resp.UserId)
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) {
		return
	}

	ctx.JSON(http.StatusOK, &LoginRes{
		StatusCode: apiutil.StatusOK,
		UserId:     resp.UserId,
		Token:      token,
	})
}

type LoginRes struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

func (api *UserAPI) Login(c context.Context, ctx *app.RequestContext) {
	params := &LoginOrRegisterReq{}
	if err := ctx.Bind(params); apiutil.HandleError(ctx, err, apiutil.ErrInvalidParams) {
		return
	}

	resp, err := api.userClient.Login(c, &user.LoginReq{
		Username: params.Username,
		Password: params.Password,
	})
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
		return
	}

	token, err := jwtutil.GenerateTokenWithUserId(resp.UserId)
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) {
		return
	}

	ctx.JSON(http.StatusOK, &LoginRes{
		StatusCode: apiutil.StatusOK,
		UserId:     resp.UserId,
		Token:      token,
	})
}

type UserInfoReq struct {
	ToUserId int64 `query:"user_id,required"`
}

type UserInfoRes struct {
	StatusCode int            `json:"status_code"`
	StatusMsg  string         `json:"status_msg"`
	Info       *user.UserInfo `json:"user"`
}

func (api *UserAPI) UserInfo(c context.Context, ctx *app.RequestContext) {
	params := &UserInfoReq{}
	if err := ctx.Bind(params); apiutil.HandleError(ctx, err, apiutil.ErrInvalidParams) {
		return
	}

	resp, err := api.userClient.GetUserInfo(c, &user.GetUserInfoReq{
		UserId:   ctx.GetInt64(middleware.CTX_USER_ID),
		ToUserId: params.ToUserId,
	})
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
		return
	}

	ctx.JSON(http.StatusOK, &UserInfoRes{
		StatusCode: apiutil.StatusOK,
		Info:       resp.User,
	})
}
