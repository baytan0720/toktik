package feed

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
	"toktik/internal/video/kitex_gen/video"
	"toktik/internal/video/kitex_gen/video/videoservice"
)

type FeedApi struct {
	videoClient videoservice.Client
}

func NewFeedApi(r discovery.Resolver) *FeedApi {
	return &FeedApi{
		videoClient: videoservice.MustNewClient("video", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}

func (api *FeedApi) Routes() []apiutil.Route {
	return []apiutil.Route{
		{
			Method:  http.MethodGet,
			Path:    "/douyin/feed",
			Handler: api.Feed,
			Hooks:   []app.HandlerFunc{middleware.SoftAuthCheck},
		},
	}
}

type FeedResp struct {
	StatusCode int                `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	NextTime   int64              `json:"next_time"`
	VideoList  []*video.VideoInfo `json:"video_list"`
}

func (api *FeedApi) Feed(c context.Context, ctx *app.RequestContext) {
	latestTime, err := strconv.ParseInt(ctx.Query("latest_time"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, &FeedResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  "invalid latest_time",
		})
		return
	}

	resp, err := api.videoClient.Feed(c, &video.FeedReq{
		UserId:     ctx.GetInt64(middleware.CTX_USER_ID),
		LatestTime: latestTime,
	})

	if err != nil {
		ctx.JSON(http.StatusOK, &FeedResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusOK, &FeedResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  resp.ErrMsg,
		})
		return
	}

	ctx.JSON(http.StatusOK, &FeedResp{
		StatusCode: apiutil.StatusOK,
		VideoList:  resp.VideoList,
		NextTime:   resp.NextTime,
	})

}
