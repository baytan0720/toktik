package comment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
)

type CommentApi struct{}

func NewCommentApi() *CommentApi {
	return &CommentApi{}
}

func (api *CommentApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/comment/action",
			Handler: api.Action,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/comment/list",
			Handler: api.List,
		},
	}
}

func (api *CommentApi) Action(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("Action...")
}

func (api *CommentApi) List(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("List...")
}
