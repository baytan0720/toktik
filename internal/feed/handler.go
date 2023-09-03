package main

import (
	"context"
	"time"

	"toktik/internal/feed/kitex_gen/feed"
	"toktik/internal/feed/pkg/ctx"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewFeedServiceImpl(svcCtx *ctx.ServiceContext) *FeedServiceImpl {
	return &FeedServiceImpl{
		svcCtx: svcCtx,
	}
}

// Feed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) Feed(ctx context.Context, req *feed.FeedReq) (resp *feed.FeedRes, err error) {
	// TODO: Your code here...
	resp = &feed.FeedRes{}
	latestTime := req.LatestTime
	userId := req.UserId

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel() // 确保在函数结束时取消超时

	feedList, err := s.svcCtx.FeedService.GetFeed(userId, latestTime)
	if err != nil {
		resp.Status = feed.Status_ERROR
		resp.ErrMsg = err.Error()
		return resp, err
	}

	resp.VideoList = make([]*feed.VideoInfo, len(feedList))

	// 检查用户是否已登录
	if userId != -1 {
		// 获取用户点赞的视频id
		for i, video := range feedList {
			videoId := video.Id
			isFavorite, err := s.svcCtx.FeedService.CheckIsFavorite(userId, videoId)
			if err != nil {
				resp.Status = feed.Status_ERROR
				resp.ErrMsg = err.Error()
				return resp, err
			}
			resp.VideoList[i] = &feed.VideoInfo{
				IsFavorite: isFavorite,
			}
		}
	}
	return resp, nil
}
