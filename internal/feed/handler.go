package main

import (
	"context"

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
	return
}
