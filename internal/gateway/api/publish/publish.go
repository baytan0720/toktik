package publish

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"

	"toktik/internal/gateway/middleware"
	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/gateway/pkg/jwtutil"
	"toktik/internal/video/kitex_gen/video"
	"toktik/internal/video/kitex_gen/video/videoservice"
)

type PublishApi struct {
	publishClient videoservice.Client
}

func NewPublishApi(r discovery.Resolver) *PublishApi {
	return &PublishApi{
		publishClient: videoservice.MustNewClient("video", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}

func (api *PublishApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/douyin/publish/action/",
			Handler: api.Action,
		},
		{
			Method:  http.MethodGet,
			Path:    "/douyin/publish/list/",
			Handler: api.List,
		},
	}
}

type PublishReq struct {
	Token string                `form:"token"` // 用户鉴权token
	Data  *multipart.FileHeader `form:"data"`  // 视频数据
	Title string                `form:"title"` // 视频标题
}

type PublishRes struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func (api *PublishApi) Action(c context.Context, ctx *app.RequestContext) {
	body := &PublishReq{}

	if err := ctx.Bind(&body); apiutil.HandleError(ctx, err, apiutil.ErrInvalidParams) {
		return
	}

	// The token is put into the form data in this request,so we authenticate here
	claim, err := jwtutil.ParseToken(body.Token)
	if apiutil.HandleError(ctx, err, err) {
		return
	}

	file, err := body.Data.Open()
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) {
		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) {
		return
	}

	if resp, err := api.publishClient.PublishVideo(c, &video.PublishVideoReq{
		UserId: claim.UserId,
		Data:   fileBytes,
		Title:  body.Title,
	}); apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
		return
	}

	ctx.JSON(http.StatusOK, &PublishRes{
		StatusCode: apiutil.StatusOK,
	})
}

type ListReq struct {
	UserId int64 `query:"user_id"`
}

type ListRes struct {
	StatusCode int                `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	VideoList  []*video.VideoInfo `json:"video_list"`
}

func (api *PublishApi) List(c context.Context, ctx *app.RequestContext) {
	params := &ListReq{}
	if err := ctx.Bind(params); apiutil.HandleError(ctx, err, apiutil.ErrInvalidParams) {
		return
	}

	resp, err := api.publishClient.ListVideo(c, &video.ListVideoReq{
		UserId:   ctx.GetInt64(middleware.CTX_USER_ID),
		ToUserId: params.UserId,
	})
	if apiutil.HandleError(ctx, err, apiutil.ErrInternalError) || apiutil.HandleRpcError(ctx, int32(resp.Status), resp.ErrMsg) {
		return
	}

	ctx.JSON(http.StatusOK, &ListRes{
		StatusCode: apiutil.StatusOK,
		VideoList:  resp.VideoList,
	})
}
