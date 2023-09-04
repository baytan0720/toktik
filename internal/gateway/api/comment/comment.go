package comment

import (
	"context"
	"net/http"
	"strconv"
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
			Path:    "/douyin/comment/action",
			Handler: api.Action,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/comment/list",
			Handler: api.List,
			Hooks:   []app.HandlerFunc{middleware.AuthCheck},
		},
	}
}

type ActionResp struct {
	StatusCode int                  `json:"status_code"`
	StatusMsg  string               `json:"status_msg"`
	Comment    *comment.CommentInfo `json:"comment"`
}

func (api *CommentApi) Action(c context.Context, ctx *app.RequestContext) {
	actionType, err1 := strconv.Atoi(ctx.Query("action_type"))
	videoId, err2 := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err1 != nil || err2 != nil {
		var err error
		if err1 != nil {
			err = err1
		} else {
			err = err2
		}
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}
	switch actionType {
	case 1:
		resp, err := api.commentClient.CreateComment(c, &comment.CreateCommentReq{
			UserId:  ctx.GetInt64(middleware.CTX_USER_ID),
			VideoId: videoId,
			Content: ctx.Query("comment_text"),
		})
		if err != nil {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		} else if resp.Status != 0 {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  resp.ErrMsg,
			})
			return
		}
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusOK,
			Comment:    resp.Comment,
		})
		break
	case 2:
		commentId, err := strconv.ParseInt(ctx.Query("comment_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		}
		resp, err := api.commentClient.DeleteComment(c, &comment.DeleteCommentReq{
			UserId:    ctx.GetInt64(middleware.CTX_USER_ID),
			CommentId: commentId,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		} else if resp.Status != 0 {
			ctx.JSON(http.StatusOK, &ActionResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  resp.ErrMsg,
			})
			return
		}
		ctx.JSON(http.StatusOK, &ActionResp{
			StatusCode: apiutil.StatusOK,
		})
		break
	}
}

type ListResp struct {
	StatusCode  int
	StatusMsg   string
	CommentList []*comment.CommentInfo
}

func (api *CommentApi) List(c context.Context, ctx *app.RequestContext) {
	userId := ctx.GetInt64(middleware.CTX_USER_ID)
	videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, &ListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}
	resp, err := api.commentClient.ListComment(c, &comment.ListCommentReq{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &ListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &ListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}

	ctx.JSON(http.StatusOK, &ListResp{
		StatusCode:  apiutil.StatusOK,
		CommentList: resp.CommentList,
	})
}
