package main

import (
	"context"
	"log"
	"os"
	"toktik/internal/video/kitex_gen/video"
	"toktik/internal/video/pkg/ctx"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewVideoServiceImpl(svcCtx *ctx.ServiceContext) *VideoServiceImpl {
	return &VideoServiceImpl{
		svcCtx: svcCtx,
	}
}

// ListVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) ListVideo(ctx context.Context, req *video.ListVideoReq) (resp *video.ListVideoRes, err error) {
	resp = &video.ListVideoRes{}
	videoList, err := s.svcCtx.VideoService.ListVideo(req.UserId)
	if err != nil {
		return nil, err
	}
	resp.video_list = make([]*video.Video, len(videoList))
	for i, video := range videoList {
		if res, err := s.svcCtx.VideoClient.GetVideo(ctx, &video.GetVideoPeq{}); err == nil {
			resp.video_list[i] = video
		} else {
			log.Println("get video info failed:", err)
		}
		return resp, nil
	}
	return resp, nil
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, req *video.PublishVideoReq) (resp *video.PublishVideoRes, err error) {
	// 创建一个视频上传客户端
	client := s.svcCtx.VideoClient

	// 创建一个视频文件流
	file, err := os.Open(req.PlayUrl)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 创建一个视频结构体
	video := &video.Video{
		Title:    req.Title,
		UserId: req.UserId
		PlayUrl: req.PlayUrl,
		CoverUrl:  req.CoverUrl,
	}
	// 调用视频上传客户端的上传方法，并传入视频文件流和视频结构体
	if resp, err := client.PublishVideo(ctx, video, file); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideo(ctx context.Context, req *video.GetVideoReq) (*video.GetVideoRes, error) {
	resp := &video.GetVideoRes{}
	if res, err := s.svcCtx.VideoClient.GetVideo(ctx, &video.GetVideoReq{VideoId: req.VideoId}); err != nil {
		log.Println("get video info failed")
		return nil, err
	} else {
		resp, err = s.svcCtx.VideoService.GetVideo(req.VideoId)
		return resp, err
	}
}

// GetWorkCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetWorkCount(ctx context.Context, req *video.GetWorkCountReq) (resp *video.GetWorkCountRes, err error) {
	// TODO: Your code here...
	return
}
