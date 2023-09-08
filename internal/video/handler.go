package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/minio/minio-go/v6"

	"toktik/internal/comment/kitex_gen/comment"
	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/user/kitex_gen/user"
	"toktik/internal/video/kitex_gen/video"
	"toktik/internal/video/pkg/ctx"
	"toktik/internal/video/pkg/snowflake"
	"toktik/pkg/config"
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

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, req *video.PublishVideoReq) (resp *video.PublishVideoRes, _ error) {
	resp = &video.PublishVideoRes{}
	filename := snowflake.Generate()
	videoBucket := config.Conf.GetString(config.KEY_MINIO_VIDEO_BUCKET)
	coverBucket := config.Conf.GetString(config.KEY_MINIO_COVER_BUCKET)
	minioExpose := config.Conf.GetString(config.KEY_MINIO_EXPOSE)
	mp4Type := "video/mp4"

	// 上传视频
	reader := bytes.NewReader(req.Data)
	if _, err := s.svcCtx.MinioClient.PutObject(videoBucket, filename+".mp4", reader, reader.Size(), minio.PutObjectOptions{ContentType: mp4Type}); err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	playUrl := fmt.Sprintf("http://%s/%s/%s.mp4", minioExpose, videoBucket, filename)
	coverUrl := fmt.Sprintf("http://%s/%s/%s.jpg", minioExpose, coverBucket, filename)

	// 异步生成封面
	err := s.svcCtx.MQ.Publish(filename)
	if err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	if err := s.svcCtx.VideoService.CreateVideo(req.UserId, req.Title, playUrl, coverUrl); err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	return
}

// ListVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) ListVideo(ctx context.Context, req *video.ListVideoReq) (resp *video.ListVideoRes, _ error) {
	resp = &video.ListVideoRes{}

	videos, err := s.svcCtx.VideoService.ListVideoByUserId(req.ToUserId)
	if err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	id2VideoInfo := make(map[int64]*video.VideoInfo)
	resp.VideoList = make([]*video.VideoInfo, 0, len(videos))
	videoIdList := make([]int64, 0, len(videos))
	for _, v := range videos {
		videoInfo := video.VideoInfo{
			Id:       v.Id,
			Title:    v.Title,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
		}
		resp.VideoList = append(resp.VideoList, &videoInfo)
		id2VideoInfo[v.Id] = &videoInfo
		videoIdList = append(videoIdList, v.Id)
	}

	wg := sync.WaitGroup{}

	// get user info
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
			UserId:   req.UserId,
			ToUserId: req.ToUserId,
		})
		if err != nil || res.Status != user.Status_OK {
			return
		}

		userInfo := convert2VideoUserInfo(res.User)
		for _, v := range resp.VideoList {
			v.Author = userInfo
		}
	}()

	// get favorite count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.FavoriteClient.GetVideoFavoriteInfo(ctx, &favorite.GetVideoFavoriteInfoReq{
			UserId:      req.UserId,
			VideoIdList: videoIdList,
		})
		if err != nil || res.Status != favorite.Status_OK {
			return
		}

		for _, info := range res.FavoriteInfoList {
			videoInfo := id2VideoInfo[info.VideoId]
			videoInfo.FavoriteCount = info.Count
			videoInfo.IsFavorite = info.IsFavorite
		}
	}()

	// get comment count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.CommentClient.GetCommentCount(ctx, &comment.GetCommentCountReq{
			VideoIdList: videoIdList,
		})
		if err != nil || res.Status != comment.Status_OK {
			return
		}

		for _, info := range res.CommentCountList {
			videoInfo := id2VideoInfo[info.VideoId]
			videoInfo.CommentCount = info.Count
		}
	}()

	wg.Wait()

	return
}

// GetVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideo(ctx context.Context, req *video.GetVideoReq) (resp *video.GetVideoRes, err error) {
	resp = &video.GetVideoRes{}

	videos, err := s.svcCtx.VideoService.GetVideoByIds(req.VideoId)
	if err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	id2VideoInfo := make(map[int64]*video.VideoInfo)
	resp.VideoList = make([]*video.VideoInfo, 0, len(videos))
	videoIdList := make([]int64, 0, len(videos))
	userIdList := make([]int64, 0, len(videos))
	for _, v := range videos {
		videoInfo := video.VideoInfo{
			Id:       v.Id,
			Title:    v.Title,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
		}
		resp.VideoList = append(resp.VideoList, &videoInfo)
		id2VideoInfo[v.Id] = &videoInfo
		videoIdList = append(videoIdList, v.Id)
		userIdList = append(userIdList, v.UserId)
	}

	wg := sync.WaitGroup{}

	// get user info
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.UserClient.GetUserInfos(ctx, &user.GetUserInfosReq{
			UserId:    req.UserId,
			ToUserIds: userIdList,
		})
		if err != nil || res.Status != user.Status_OK {
			return
		}

		id2UserInfo := make(map[int64]*video.UserInfo)
		for _, userInfo := range res.Users {
			id2UserInfo[userInfo.Id] = convert2VideoUserInfo(userInfo)
		}

		for _, info := range resp.VideoList {
			info.Author = id2UserInfo[info.Author.Id]
		}
	}()

	// get favorite count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.FavoriteClient.GetVideoFavoriteInfo(ctx, &favorite.GetVideoFavoriteInfoReq{
			UserId:      req.UserId,
			VideoIdList: videoIdList,
		})
		if err != nil || res.Status != favorite.Status_OK {
			return
		}

		for _, info := range res.FavoriteInfoList {
			videoInfo := id2VideoInfo[info.VideoId]
			videoInfo.FavoriteCount = info.Count
			videoInfo.IsFavorite = info.IsFavorite
		}
	}()

	// get comment count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.CommentClient.GetCommentCount(ctx, &comment.GetCommentCountReq{
			VideoIdList: videoIdList,
		})
		if err != nil || res.Status != comment.Status_OK {
			return
		}

		for _, info := range res.CommentCountList {
			videoInfo := id2VideoInfo[info.VideoId]
			videoInfo.CommentCount = info.Count
		}
	}()

	wg.Wait()

	return
}

// GetWorkCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetWorkCount(ctx context.Context, req *video.GetWorkCountReq) (resp *video.GetWorkCountRes, _ error) {
	resp = &video.GetWorkCountRes{}

	for _, userId := range req.UserIdList {
		count, err := s.svcCtx.VideoService.CountWork(userId)
		if err != nil {
			resp.Status = video.Status_ERROR
			resp.ErrMsg = err.Error()
			return
		}
		resp.WorkCountList = append(resp.WorkCountList, &video.WorkCountInfo{
			UserId:    userId,
			WorkCount: count,
		})
	}

	return
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedReq) (resp *video.FeedRes, err error) {
	resp = &video.FeedRes{}

	videoList, err := s.svcCtx.VideoService.GetFeed(req.LatestTime)
	if err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}
	if len(videoList) == 0 {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = "没有更多视频了"
		return
	}
	resp.NextTime = videoList[len(videoList)-1].CreatedAt.UnixMilli()

	id2VideoInfo := make(map[int64]*video.VideoInfo)
	userIdList := make([]int64, 0, len(videoList))
	resp.VideoList = make([]*video.VideoInfo, 0, len(videoList))
	videoIdList := make([]int64, 0, len(videoList))
	for _, v := range videoList {
		videoInfo := video.VideoInfo{
			Id:       v.Id,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		}
		resp.VideoList = append(resp.VideoList, &videoInfo)
		id2VideoInfo[v.Id] = &videoInfo
		videoIdList = append(videoIdList, v.Id)
		userIdList = append(userIdList, v.UserId)
	}

	wg := sync.WaitGroup{}

	// get user info
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.UserClient.GetUserInfos(ctx, &user.GetUserInfosReq{
			UserId:    req.UserId,
			ToUserIds: userIdList,
		})
		if err != nil || res.Status != user.Status_OK {
			return
		}

		id2UserInfo := make(map[int64]*video.UserInfo)
		for _, userInfo := range res.Users {
			id2UserInfo[userInfo.Id] = convert2VideoUserInfo(userInfo)
		}

		for _, info := range resp.VideoList {
			info.Author = id2UserInfo[info.Author.Id]
		}
	}()

	// get favorite count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.FavoriteClient.GetVideoFavoriteInfo(ctx, &favorite.GetVideoFavoriteInfoReq{
			UserId:      req.UserId,
			VideoIdList: videoIdList,
		})
		if err != nil || res.Status != favorite.Status_OK {
			return
		}

		for _, info := range res.FavoriteInfoList {
			videoInfo := id2VideoInfo[info.VideoId]
			videoInfo.FavoriteCount = info.Count
			videoInfo.IsFavorite = info.IsFavorite
		}
	}()

	// get comment count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.CommentClient.GetCommentCount(ctx, &comment.GetCommentCountReq{
			VideoIdList: videoIdList,
		})
		if err != nil || res.Status != comment.Status_OK {
			return
		}

		for _, info := range res.CommentCountList {
			videoInfo := id2VideoInfo[info.VideoId]
			videoInfo.CommentCount = info.Count
		}
	}()

	wg.Wait()

	return resp, nil
}

func convert2VideoUserInfo(user *user.UserInfo) *video.UserInfo {
	return &video.UserInfo{
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
