package main

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"toktik/internal/video/kitex_gen/video"

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
	resp = &favorite.FavoriteRes{}
	videoId := []int64{req.VideoId}
	if res, err := s.svcCtx.VideoClient.GetVideo(ctx, &video.GetVideoReq{
		VideoId: videoId,
	}); err != nil {
		log.Println("get video info failed", err)
		return nil, err
	} else if res.Status != video.Status_OK {
		log.Println("video not exit")
		resp.Status = favorite.Status_ERROR
		resp.ErrMsg = "video not exist"
	} else {
		err = s.svcCtx.FavoriteService.Favorite(req.VideoId, req.UserId)
		if err != nil {
			// internal error
			return nil, err
		}
		return resp, nil
	}
	return
}

// UnFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) UnFavorite(ctx context.Context, req *favorite.UnFavoriteReq) (resp *favorite.UnFavoriteRes, err error) {
	resp = &favorite.UnFavoriteRes{}
	err = s.svcCtx.FavoriteService.UnFavorite(req.VideoId, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Status = favorite.Status_ERROR
			resp.ErrMsg = err.Error()
		} else {
			return nil, err
		}
	}
	return resp, nil
}

// ListFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) ListFavorite(ctx context.Context, req *favorite.ListFavoriteReq) (resp *favorite.ListFavoriteRes, err error) {
	resp = &favorite.ListFavoriteRes{}
	videoIdList, err := s.svcCtx.FavoriteService.ListFavorite(req.UserId)
	if err != nil {
		return nil, err
	}
	resp.VideoList = make([]*favorite.VideoInfo, len(videoIdList))

	if res, err := s.svcCtx.VideoClient.GetVideo(ctx, &video.GetVideoReq{
		VideoId: videoIdList,
	}); err != nil {
		log.Println("get video info failed", err)
	} else {
		for i, videoId := range videoIdList {
			favoriteCount, err := s.svcCtx.FavoriteService.CountFavorite(videoId)
			if err != nil {
				return nil, err
			}
			resp.VideoList[i] = convert2FavoriteVideoInfo(res.VideoList[i])
			resp.VideoList[i].FavoriteCount = favoriteCount
		}
	}

	return resp, err
}

// GetFavoriteInfo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteInfo(ctx context.Context, req *favorite.GetFavoriteInfoReq) (resp *favorite.GetFavoriteInfoRes, err error) {
	resp = &favorite.GetFavoriteInfoRes{}
	resp.FavoriteInfoList = make([]*favorite.FavoriteInfo, 0, len(req.VideoIdList))

	for _, videoId := range req.VideoIdList {
		favoriteCount, err := s.svcCtx.FavoriteService.CountFavorite(videoId)
		if err != nil {
			return nil, err
		}
		isFavorite, err := s.svcCtx.FavoriteService.IsFavorite(videoId, req.UserId)
		if err != nil {
			return nil, err
		}
		resp.FavoriteInfoList = append(resp.FavoriteInfoList, &favorite.FavoriteInfo{
			VideoId:    videoId,
			Count:      favoriteCount,
			IsFavorite: isFavorite,
		})
	}
	return resp, err
}

func convert2FavoriteVideoInfo(video *video.Video) *favorite.VideoInfo {

	return &favorite.VideoInfo{
		Id:           video.Id,
		Author:       convert2FavoriteUserInfo(video.Author),
		PlayUrl:      video.PlayUrl,
		CoverUrl:     video.CoverUrl,
		CommentCount: video.CommentCount,
		IsFavorite:   true,
		Title:        video.Title,
	}
}

func convert2FavoriteUserInfo(user *video.User) *favorite.UserInfo {
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
