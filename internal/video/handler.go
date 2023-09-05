package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"sync"
	"time"

	"github.com/minio/minio-go/v6"
	ffmpeg "github.com/u2takey/ffmpeg-go"

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
	videoBucket := config.Conf.Get(config.KEY_MINIO_VIDEO_BUCKET).(string)
	coverBucket := config.Conf.Get(config.KEY_MINIO_COVER_BUCKET).(string)
	minioEndpoint := config.Conf.Get(config.KEY_MINIO_ENDPOINT).(string)
	mp4Type := "video/mp4"
	JpegType := "image/jpeg"

	// 上传视频
	reader := bytes.NewReader(req.Data)
	if _, err := s.svcCtx.MinioClient.PutObject(videoBucket, filename+".mp4", reader, reader.Size(), minio.PutObjectOptions{ContentType: mp4Type}); err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	playUrl := fmt.Sprintf("http://%s/%s/%s.mp4" + minioEndpoint + videoBucket + filename)
	coverUrl := fmt.Sprintf("http://%s/%s/%s.jpg" + minioEndpoint + coverBucket + filename)

	// 异步生成封面
	go func() {
		// 获取封面
		coverData, _ := readFrameAsJpeg(playUrl)

		//上传封面
		coverReader := bytes.NewReader(coverData)
		_, _ = s.svcCtx.MinioClient.PutObject(coverBucket, filename+".jpg", coverReader, coverReader.Size(), minio.PutObjectOptions{ContentType: JpegType})
	}()

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

// ReadFrameAsJpeg
// 从视频流中截取一帧并返回 需要在本地环境中安装ffmpeg并将bin添加到环境变量
func readFrameAsJpeg(filePath string) ([]byte, error) {
	reader := bytes.NewBuffer(nil)
	err := ffmpeg.Input(filePath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(reader, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, img, nil)

	return buf.Bytes(), err
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

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedReq) (resp *video.FeedRes, err error) {
	// TODO: Your code here...
	resp = &video.FeedRes{}

	latestTime := req.LatestTime
	userId := req.UserId

	if latestTime == "" {
		nowTime := time.Now().UnixMilli()
		curTime := time.Unix(0, nowTime*int64(time.Millisecond))
		latestTime = curTime.Format("2006-01-02 15:04:05")
	}

	videoList, nextTime, err := s.svcCtx.VideoService.GetFeed(userId, latestTime)

	id2VideoInfo := make(map[int64]*video.VideoInfo)
	userIdList := make([]int64, 0, len(videoList))
	resp.VideoList = make([]*video.VideoInfo, 0, len(videoList))
	videoIdList := make([]int64, 0, len(videoList))

	if err != nil {
		resp.Status = video.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	for _, v := range videoList {
		videoInfo := video.VideoInfo{
			Id:       v.Id,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		}
		resp.VideoList = append(resp.VideoList, &videoInfo)
		userIdList = append(userIdList, v.UserId)
	}

	resp.NextTime = nextTime

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
