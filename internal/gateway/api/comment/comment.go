package comment

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"

	"toktik/internal/comment/kitex_gen/comment"
	"toktik/internal/comment/kitex_gen/comment/commentservice"
	"toktik/internal/gateway/middleware"
	"toktik/internal/gateway/pkg/apiutil"
)

type CommentApi struct {
	commentClient commentservice.Client
}

func NewCommentApi(r discovery.Resolver) *CommentApi {
	return &CommentApi{
		commentClient: commentservice.MustNewClient("comment", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}

func (api *CommentApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/comment/action/",
			Handler: api.Action,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/comment/list/",
			Handler: api.List,
			Hooks:   []app.HandlerFunc{middleware.SoftAuthCheck},
		},
	}
}

type ActionReq struct {
	ActionType  int    `query:"action_type,required"`
	VideoId     int64  `query:"video_id"`
	CommentText string `query:"comment_text"`
	CommentId   int64  `query:"comment_id"`
}

type ActionRes struct {
	StatusCode int                  `json:"status_code"`
	StatusMsg  string               `json:"status_msg"`
	Comment    *comment.CommentInfo `json:"comment"`
}

func (api *CommentApi) Action(c context.Context, ctx *app.RequestContext) {
	params := &ActionReq{}
	if err := ctx.Bind(params); apiutil.HandleError(ctx, err, apiutil.ErrInvalidParams) {
		return
	}

	switch params.ActionType {
	case 1:
		resp, err := api.commentClient.CreateComment(c, &comment.CreateCommentReq{
			UserId:  ctx.GetInt64(middleware.CTX_USER_ID),
			VideoId: params.VideoId,
			Content: params.CommentText,
		})
		if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
			return
		}
		ctx.JSON(http.StatusOK, &ActionRes{
			StatusCode: apiutil.StatusOK,
			Comment:    resp.Comment,
		})
	case 2:
		resp, err := api.commentClient.DeleteComment(c, &comment.DeleteCommentReq{
			UserId:    ctx.GetInt64(middleware.CTX_USER_ID),
			CommentId: params.CommentId,
		})
		if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
			return
		}
		ctx.JSON(http.StatusOK, &ActionRes{
			StatusCode: apiutil.StatusOK,
		})
	default:
		ctx.JSON(http.StatusOK, &ActionRes{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  apiutil.ErrInvalidParams.Error(),
		})
	}
}

type ListReq struct {
	VideoId int64 `query:"video_id,required"`
}

type ListRes struct {
	StatusCode  int                    `json:"status_code"`
	StatusMsg   string                 `json:"status_msg"`
	CommentList []*comment.CommentInfo `json:"comment_list"`
}

func (api *CommentApi) List(c context.Context, ctx *app.RequestContext) {
	params := &ListReq{}
	if err := ctx.Bind(params); apiutil.HandleError(ctx, err, apiutil.ErrInvalidParams) {
		return
	}

	resp, err := api.commentClient.ListComment(c, &comment.ListCommentReq{
		UserId:  ctx.GetInt64(middleware.CTX_USER_ID),
		VideoId: params.VideoId,
	})
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
		return
	}

	ctx.JSON(http.StatusOK, &ListRes{
		StatusCode:  apiutil.StatusOK,
		CommentList: resp.CommentList,
	})
}
