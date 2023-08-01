package user

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
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
	}
}

func (api *UserAPI) Register(c context.Context, ctx *app.RequestContext) {
	// TODO
}

func (api *UserAPI) Login(c context.Context, ctx *app.RequestContext) {
	// TODO
}
