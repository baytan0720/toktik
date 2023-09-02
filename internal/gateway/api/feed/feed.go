package feed

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"

	"toktik/internal/feed/kitex_gen/feed"
	"toktik/internal/feed/kitex_gen/feed/feedservice"
	"toktik/internal/gateway/middleware"
	"toktik/internal/gateway/pkg/apiutil"
)

var (
	FeedFail = "fail to get feed"
)

type FeedApi struct {
	feedClient feedservice.Client
}

func NewFeedApi(r discovery.Resolver) *FeedApi {
	return &FeedApi{
		feedClient: feedservice.MustNewClient("feed", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
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
	StatusCode int               `json:"status_code"`
	StatusMsg  string            `json:"status_msg"`
	NextTime   string            `json:"next_time"`
	VideoList  []*feed.VideoInfo `json:"video_list"`
}

func (api *FeedApi) Feed(c context.Context, ctx *app.RequestContext) {
	latestTime := ctx.Query("latest_time")

	resp, err := api.feedClient.Feed(c, &feed.FeedReq{
		UserId:     ctx.GetInt64(middleware.CTX_USER_ID),
		LatestTime: latestTime,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &FeedResp{
			StatusCode: apiutil.StatusFailed,
			StatusMsg:  err.Error(),
		})
		return
	} else if resp.Status != 0 {
		ctx.JSON(http.StatusBadRequest, &FeedResp{
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
