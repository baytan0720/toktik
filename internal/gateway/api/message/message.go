package message

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
	"toktik/internal/message/kitex_gen/message"
	"toktik/internal/message/kitex_gen/message/messageservice"
)

type MessageApi struct {
	messageClient messageservice.Client
}

func NewMessageApi(r discovery.Resolver) *MessageApi {
	return &MessageApi{
		messageClient: messageservice.MustNewClient("message", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}

func (api *MessageApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/message/action",
			Handler: api.Action,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/message/chat",
			Handler: api.Chat,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
	}
}

type ChatResp struct {
	StatusCode  int                    `json:"status_code"`
	StatusMsg   string                 `json:"status_msg"`
	MessageList []*message.MessageInfo `json:"message_list"`
}

func (api *MessageApi) Chat(c context.Context, ctx *app.RequestContext) {
	userId := ctx.GetInt64(middleware.CTX_USER_ID)
	toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &ChatResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}
	resp, err := api.messageClient.ListMessage(c, &message.ListMessageReq{
		UserId:   userId,
		ToUserId: toUserId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &ChatResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusBadRequest, &ChatResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, &ChatResp{
		StatusCode:  apiutil.StatusOK,
		MessageList: resp.MessageList,
	})
}

type ActionResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func (api *MessageApi) Action(c context.Context, ctx *app.RequestContext) {
	actionType, err1 := strconv.Atoi(ctx.Query("action_type"))
	toUserId, err2 := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err1 != nil || err2 != nil {
		var err error
		if err1 != nil {
			err = err1
		} else {
			err = err2
		}
		ctx.JSON(http.StatusBadRequest, &ActionResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}
	switch actionType {
	case 1:
		resp, err := api.messageClient.SendMessage(c, &message.SendMessageReq{
			UserId:   ctx.GetString(middleware.CTX_USER_ID),
			ToUserId: toUserId,
			Content:  ctx.Query("content"),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		} else if resp.Status != 0 {
			ctx.JSON(http.StatusBadRequest, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  resp.ErrMsg,
			})
			return
		}
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusOK,
		})
		break
	default:
		ctx.JSON(http.StatusBadRequest, &ActionResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  "fail to send message",
		})
		return
	}
}
