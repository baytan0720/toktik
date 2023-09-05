package publish

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
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

var (
	PublishFail = "fail to publish"
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

type ListResp struct {
	StatusCode int                `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	VideoList  []*video.VideoInfo `json:"video_list"`
}

func (api *PublishApi) List(c context.Context, ctx *app.RequestContext) {
	userId := ctx.GetInt64(middleware.CTX_USER_ID)
	toUserId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, &ListResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}
	resp, err := api.publishClient.ListVideo(c, &video.ListVideoReq{
		UserId:   userId,
		ToUserId: toUserId,
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
		StatusCode: apiutil.StatusOK,
		VideoList:  resp.VideoList,
	})
}

type PublishResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type PublishReq struct {
	Token string                `json:"token" form:"token" query:"token"` // 用户鉴权token
	Data  *multipart.FileHeader `form:"data"`                                               // 视频数据
	Title string                `json:"title" form:"title" query:"title"` // 视频标题
}

func (api *PublishApi) Action(c context.Context, ctx *app.RequestContext) {
	body := &PublishReq{}

	err := ctx.Bind(&body)
	if err != nil {
		ctx.JSON(http.StatusOK, &PublishResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}

	//The token is put into the body in this request,so we authenticate here
	j := jwtutil.NewJwtUtil()
	claim, err := j.ParseToken(body.Token)
	if err != nil {
		ctx.JSON(http.StatusOK, &PublishResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}

	file, err := body.Data.Open()
	if err != nil {
		ctx.JSON(http.StatusOK, &PublishResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			ctx.JSON(http.StatusOK, &PublishResp{
				StatusCode: apiutil.StatusFailed,
				StatusMsg:  err.Error(),
			})
			return
		}
	}(file)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusOK, &PublishResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	}

	resp, err := api.publishClient.PublishVideo(c, &video.PublishVideoReq{
		UserId: claim.UserId,
		Data:   fileBytes,
		Title:  body.Title,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, &PublishResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &PublishResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}

	ctx.JSON(http.StatusOK, &PublishResp{
		StatusCode: apiutil.StatusOK,
	})
}
