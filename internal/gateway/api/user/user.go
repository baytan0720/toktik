package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"

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
			Path:    "/douyin/user/register",
			Handler: api.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/douyin/user/login",
			Handler: api.Login,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/user",
			Handler: api.UserInfo,
		},
	}
}

func (api *UserAPI) UserInfo(c context.Context, ctx *app.RequestContext) {
	// TODO
	fmt.Println("UserInfo")
}

func (api *UserAPI) Register(c context.Context, ctx *app.RequestContext) {
	// TODO
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
		ctx.JSON(http.StatusInternalServerError, &LoginResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusBadRequest, &LoginResp{
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
