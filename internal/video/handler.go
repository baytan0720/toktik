package main

import (
	"context"

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
	// TODO: Your code here...
	return
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, req *video.PublishVideoReq) (resp *video.PublishVideoRes, err error) {
	// TODO: Your code here...
	return
}

// GetVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideo(ctx context.Context, req *video.GetVideoReq) (resp *video.GetVideoRes, err error) {
	// TODO: Your code here...
	return
}

// GetWorkCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetWorkCount(ctx context.Context, req *video.GetWorkCountReq) (resp *video.GetWorkCountRes, err error) {
	// TODO: Your code here...
	return
}
