package publish

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
)

type PublishApi struct{}

func NewPublishApi() *PublishApi {
	return &PublishApi{}
}

func (api *PublishApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/publish/action",
			Handler: api.Action,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/publish/list",
			Handler: api.List,
		},
	}
}

func (api *PublishApi) List(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("list")
}

func (api *PublishApi) Action(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("action")
}
