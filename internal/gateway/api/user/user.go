package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/gateway/pkg/jwtutil"
)

type UserAPI struct{}

func NewUserAPI() *UserAPI {
	return &UserAPI{}
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

func (api *UserAPI) Login(c context.Context, ctx *app.RequestContext) {
	// TODO

	//登陆成功后
	userid := 111
	j := jwtutil.NewJwtUtil()
	token, _ := j.GenerateToken(jwtutil.CreateClaims(userid))
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
