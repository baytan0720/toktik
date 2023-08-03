package favorite

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
)

type FavoriteApi struct{}

func NewFavoriteApi() *FavoriteApi {
	return &FavoriteApi{}
}

func (api *FavoriteApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/favorite/action",
			Handler: api.Action,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/favorite/list",
			Handler: api.List,
		},
	}
}

func (api *FavoriteApi) Action(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("Action...")
}

func (api *FavoriteApi) List(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("List...")
}
