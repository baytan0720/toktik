package main

import (
	"context"

	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/favorite/pkg/ctx"
	"toktik/internal/video/kitex_gen/video"
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

var (
	errFavoriteFailed     = "点赞失败"
	errVideoNotFound      = "视频不存在"
	errUnFavoriteFailed   = "取消点赞失败"
	errListFavoriteFailed = "获取点赞列表失败"
)

// Favorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) Favorite(ctx context.Context, req *favorite.FavoriteReq) (resp *favorite.FavoriteRes, _ error) {
	resp = &favorite.FavoriteRes{}

	res, err := s.svcCtx.VideoClient.IsExist(ctx, &video.IsExistReq{
		VideoId: req.VideoId,
	})
	if err != nil || res.Status != video.Status_OK {
		resp.Status = favorite.Status_ERROR
		resp.ErrMsg = errFavoriteFailed
		return
	}
	if !res.IsExist {
		resp.Status = favorite.Status_ERROR
		resp.ErrMsg = errVideoNotFound
		return
	}

	err = s.svcCtx.FavoriteService.Favorite(req.VideoId, req.UserId)
	if err != nil {
		resp.Status = favorite.Status_ERROR
		resp.ErrMsg = errFavoriteFailed
		return
	}

	return
}

// UnFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) UnFavorite(ctx context.Context, req *favorite.UnFavoriteReq) (resp *favorite.UnFavoriteRes, _ error) {
	resp = &favorite.UnFavoriteRes{}

	err := s.svcCtx.FavoriteService.UnFavorite(req.VideoId, req.UserId)
	if err != nil {
		resp.Status = favorite.Status_ERROR
		resp.ErrMsg = errUnFavoriteFailed
		return
	}

	return
}

// ListFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) ListFavorite(ctx context.Context, req *favorite.ListFavoriteReq) (resp *favorite.ListFavoriteRes, _ error) {
	resp = &favorite.ListFavoriteRes{}
	videoIdList, err := s.svcCtx.FavoriteService.ListFavoriteByUserId(req.UserId)
	if err != nil {
		resp.Status = favorite.Status_ERROR
		resp.ErrMsg = errListFavoriteFailed
		return
	}
	resp.VideoList = make([]*favorite.VideoInfo, len(videoIdList))

	res, err := s.svcCtx.VideoClient.GetVideo(ctx, &video.GetVideoReq{
		VideoId:     videoIdList,
		AllFavorite: true,
	})
	if err != nil || res.Status != video.Status_OK {
		return
	}

	id2Video := make(map[int64]*favorite.VideoInfo, len(res.VideoList))
	for _, v := range res.VideoList {
		id2Video[v.Id] = convert2FavoriteVideoInfo(v)
	}

	for i, id := range videoIdList {
		resp.VideoList[i] = id2Video[id]
	}

	return resp, err
}

// GetVideoFavoriteInfo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetVideoFavoriteInfo(ctx context.Context, req *favorite.GetVideoFavoriteInfoReq) (resp *favorite.GetVideoFavoriteInfoRes, _ error) {
	resp = &favorite.GetVideoFavoriteInfoRes{}

	if len(req.VideoIdList) == 0 {
		return
	}

	resp.FavoriteInfoList = make([]*favorite.VideoFavoriteInfo, 0, len(req.VideoIdList))
	isFavorite := make(map[int64]bool)
	if req.UserId != 0 {
		videoIdList, _ := s.svcCtx.FavoriteService.ListFavoriteByUserId(req.UserId)
		for _, id := range videoIdList {
			isFavorite[id] = true
		}
	}

	for _, videoId := range req.VideoIdList {
		favoriteCount, _ := s.svcCtx.FavoriteService.CountVideoFavorite(videoId)
		resp.FavoriteInfoList = append(resp.FavoriteInfoList, &favorite.VideoFavoriteInfo{
			VideoId:    videoId,
			Count:      favoriteCount,
			IsFavorite: isFavorite[videoId],
		})
	}

	return
}

// GetUserFavoriteInfo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetUserFavoriteInfo(ctx context.Context, req *favorite.GetUserFavoriteInfoReq) (resp *favorite.GetUserFavoriteInfoRes, err error) {
	resp = &favorite.GetUserFavoriteInfoRes{}

	resp.FavoriteInfoList = make([]*favorite.UserFavoriteInfo, 0, len(req.UserIdList))

	for _, userId := range req.UserIdList {
		favoriteCount, _ := s.svcCtx.FavoriteService.CountUserFavorite(userId)
		var totalFavorite int64
		res, err := s.svcCtx.VideoClient.ListVideoId(ctx, &video.ListVideoIdReq{
			UserId: userId,
		})
		if err == nil && res.Status == video.Status_OK {
			for _, id := range res.VideoIdList {
				count, _ := s.svcCtx.FavoriteService.CountVideoFavorite(id)
				totalFavorite += count
			}
		}

		resp.FavoriteInfoList = append(resp.FavoriteInfoList, &favorite.UserFavoriteInfo{
			UserId:         userId,
			FavoriteCount:  favoriteCount,
			TotalFavorited: totalFavorite,
		})
	}

	return
}

func convert2FavoriteVideoInfo(video *video.VideoInfo) *favorite.VideoInfo {
	return &favorite.VideoInfo{
		Id:           video.Id,
		Author:       convert2FavoriteUserInfo(video.Author),
		PlayUrl:      video.PlayUrl,
		CoverUrl:     video.CoverUrl,
		CommentCount: video.CommentCount,
		IsFavorite:   video.IsFavorite,
		Title:        video.Title,
	}
}

func convert2FavoriteUserInfo(user *video.UserInfo) *favorite.UserInfo {
	return &favorite.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
