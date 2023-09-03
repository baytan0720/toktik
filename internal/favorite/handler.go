package main

import (
	"context"

	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/favorite/pkg/ctx"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewFavoriteServiceImpl(svcCtx *ctx.ServiceContext) *FavoriteServiceImpl {
	return &FavoriteServiceImpl{
		svcCtx: svcCtx,
	}
}

// Favorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) Favorite(ctx context.Context, req *favorite.FavoriteReq) (resp *favorite.FavoriteRes, err error) {
	// TODO: Your code here...
	return
}

// UnFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) UnFavorite(ctx context.Context, req *favorite.UnFavoriteReq) (resp *favorite.UnFavoriteRes, err error) {
	// TODO: Your code here...
	return
}

// ListFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) ListFavorite(ctx context.Context, req *favorite.ListFavoriteReq) (resp *favorite.ListFavoriteRes, err error) {
	// TODO: Your code here...
	return
}

// GetVideoFavoriteInfo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetVideoFavoriteInfo(ctx context.Context, req *favorite.GetVideoFavoriteInfoReq) (resp *favorite.GetVideoFavoriteInfoRes, err error) {
	// TODO: Your code here...
	return
}

// GetUserFavoriteInfo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetUserFavoriteInfo(ctx context.Context, req *favorite.GetUserFavoriteInfoReq) (resp *favorite.GetUserFavoriteInfoRes, err error) {
	// TODO: Your code here...
	return
}
