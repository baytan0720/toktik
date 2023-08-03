package relation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
)

type RelationApi struct{}

func NewRelationApi() *RelationApi {
	return &RelationApi{}
}

func (api *RelationApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/relation/action",
			Handler: api.Action,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/relation/follow/list",
			Handler: api.FollowList,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/relation/follower/list",
			Handler: api.FollowerList,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/relation/friend/list",
			Handler: api.FriendList,
		},
	}
}

func (api *RelationApi) Action(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("Action")
}

func (api *RelationApi) FollowList(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("FollowList")
}

func (api *RelationApi) FollowerList(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("FollowerList")
}

func (api *RelationApi) FriendList(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("FriendList")
}
