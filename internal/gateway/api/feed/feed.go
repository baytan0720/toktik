package feed

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
)

type FeedApi struct{}

func NewFeedApi() *FeedApi {
	return &FeedApi{}
}

func (api *FeedApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodGet,
			Path:    "/douyin/feed",
			Handler: api.Feed,
		},
	}
}

func (api *FeedApi) Feed(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("feed")
}
