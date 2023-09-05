package user

import (
	"context"
	"net/http"
	"strconv"
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

var (
	RegisterFail = "fail to register"
	GetInfoFail  = "fail to get info"
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

type UserInfoRes struct {
	StatusCode int            `json:"status_code"`
	StatusMsg  string         `json:"status_msg"`
	Info       *user.UserInfo `json:"user"`
}

func (api *UserAPI) UserInfo(c context.Context, ctx *app.RequestContext) {
	//所请求的用户id
	toUserId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusOK, &UserInfoRes{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
			Info:       nil,
		})
		return
	}

	userId := ctx.GetInt64(middleware.CTX_USER_ID)

	resp, err := api.userClient.GetUserInfo(c, &user.GetUserInfoReq{
		UserId:   userId,
		ToUserId: toUserId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &UserInfoRes{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &UserInfoRes{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, &UserInfoRes{
		StatusCode: apiutil.StatusOK,
		Info:       resp.User,
	})
}

type RegisterRes struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

func (api *UserAPI) Register(c context.Context, ctx *app.RequestContext) {
	resp, err := api.userClient.Register(c, &user.RegisterReq{
		Username: ctx.Query("username"),
		Password: ctx.Query("password"),
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &RegisterRes{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &RegisterRes{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}
	j := jwtutil.NewJwtUtil()
	token, _ := j.GenerateToken(jwtutil.CreateClaims(resp.UserId))
	ctx.JSON(http.StatusOK, &RegisterRes{
		StatusCode: apiutil.StatusOK,
		UserId:     resp.UserId,
		Token:      token,
	})
}

type LoginResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

func (api *UserAPI) Login(c context.Context, ctx *app.RequestContext) {
	resp, err := api.userClient.Login(c, &user.LoginReq{
		Username: ctx.Query("username"),
		Password: ctx.Query("password"),
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &LoginResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &LoginResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}

	j := jwtutil.NewJwtUtil()
	token, _ := j.GenerateToken(jwtutil.CreateClaims(resp.UserId))
	ctx.JSON(http.StatusOK, &LoginResp{
		StatusCode: apiutil.StatusOK,
		UserId:     resp.UserId,
		Token:      token,
	})
}
