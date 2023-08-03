package message

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
)

type MessageApi struct{}

func NewMessageApi() *MessageApi {
	return &MessageApi{}
}

func (api *MessageApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/message/action",
			Handler: api.Action,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/message/chat",
			Handler: api.Chat,
		},
	}
}

func (api *MessageApi) Chat(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("chat")
}

func (api *MessageApi) Action(c context.Context, ctx *app.RequestContext) {
	//TODO
	fmt.Println("action")
}
