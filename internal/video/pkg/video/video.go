package video

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v7"
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type Video = model.Video

type VideoService struct {
	dbInstance    func() *gorm.DB
	redisInstance func() *goredis.Client
}

const (
	KeyId       = "Id"
	KeyUserId   = "UserId"
	KeyTitle    = "Title"
	KeyPlayUrl  = "PlayUrl"
	KeyCoverUrl = "CoverUrl"
	KeyExisted  = "Existed"
)

func NewVideoService(db func() *gorm.DB, rdb func() *goredis.Client) *VideoService {
	return &VideoService{
		dbInstance:    db,
		redisInstance: rdb,
	}
}

func (s *VideoService) CreateVideo(userId int64, title, playUrl, coverUrl string) error {
	db := s.dbInstance()
	rdb := s.redisInstance()

	video := &Video{
		UserId:   userId,
		Title:    title,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
	}

	if err := db.Create(video).Error; err != nil {
		return err
	}

	key := strconv.FormatInt(video.Id, 10)
	err := rdb.HSet(key, FormatVideoInfo(video, true)).Err()
	if err != nil {
		return err
	}
	rdb.Expire(key, generateExpireTime())

	// 将视频ID添加到Set中
	setKey := fmt.Sprintf("user_videos:%d", userId)

	_, err = rdb.SAdd(setKey, video.Id).Result()
	if err != nil {
		return err
	}
	rdb.Expire(setKey, generateExpireTime())
	return nil
}

func (s *VideoService) ListVideoByUserId(userId int64) ([]*Video, error) {
	rdb := s.redisInstance()

	setKey := fmt.Sprintf("user_videos:%d", userId)
	videoIds, err := rdb.SMembers(setKey).Result()
	if err != nil {
		return nil, err
	}

	// 未缓存，从数据库中获取
	if len(videoIds) == 0 {
		db := s.dbInstance()

		var videos []*Video
		if err := db.Where("user_id = ?", userId).Order("created_at desc").Find(&videos).Error; err != nil {
			return nil, err
		}

		// 如果数据库中没有数据，将0添加到缓存中
		if len(videos) == 0 {
			_, err := rdb.SAdd(setKey, 0).Result()
			if err != nil {
				return nil, err
			}
			rdb.Expire(setKey, generateExpireTime())
			return []*Video{}, nil
		}

		// 将视频ID添加到Set中
		for _, video := range videos {
			_, err := rdb.SAdd(setKey, video.Id).Result()
			if err != nil {
				return nil, err
			}
		}
		rdb.Expire(setKey, generateExpireTime())

		return videos, nil
	}

	// 缓存了但为空，直接返回
	if len(videoIds) == 1 && videoIds[0] == "0" {
		return []*Video{}, nil
	}

	// 命中缓存，从缓存中获取
	videos := make([]*Video, 0, len(videoIds))
	for _, key := range videoIds {
		videoData, err := rdb.HGetAll(key).Result()
		if err != nil {
			return nil, err
		}

		// 解析videoData
		video := ParseVideoInfo(videoData)
		videos = append(videos, video)
	}

	return videos, nil
}

func (s *VideoService) GetVideoByIds(videoIds []int64) ([]*Video, error) {
	rdb := s.redisInstance()

	videos := make([]*Video, 0, len(videoIds))
	for _, videoId := range videoIds {
		videoData, err := rdb.HGetAll(strconv.FormatInt(videoId, 10)).Result()
		if err != nil {
			return nil, err
		} else if len(videoData) == 0 {
			// 未命中缓存，从数据库中获取
			db := s.dbInstance()
			video := &Video{}
			if err := db.Where("id = ?", videoId).First(video).Error; err != nil {
				// 如果未找到记录，则在缓存中标识视频不存在
				if errors.Is(err, gorm.ErrRecordNotFound) {
					key := strconv.FormatInt(videoId, 10)
					err = rdb.HSet(key, FormatVideoInfo(video, false)).Err()
					if err != nil {
						return nil, err
					}
					rdb.Expire(key, generateExpireTime())
					return nil, fmt.Errorf("video not found: %d", videoId)
				}
				// internal error
				return nil, err
			} else {
				videos = append(videos, video)
			}
		} else {
			// 判断视频是否存在
			existed := videoData[KeyExisted]
			if existed != "1" {
				return nil, fmt.Errorf("video not found: %d", videoId)
			}

			// 解析videoData
			video := ParseVideoInfo(videoData)
			videos = append(videos, video)
		}
	}

	return videos, nil
}

func (s *VideoService) CountWork(userId int64) (int64, error) {
	db := s.dbInstance()

	var count int64
	if err := db.Model(&Video{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *VideoService) GetFeed(latestTime int64) ([]*Video, error) {
	db := s.dbInstance()

	videos := make([]*Video, 0)
	var timeValue time.Time

	if latestTime == 0 {
		timeValue = time.Now()
	} else {
		timeValue = time.Unix(latestTime/time.Microsecond.Nanoseconds(), 0)
	}

	if err := db.Where("created_at < ?", timeValue).Order("created_at desc").Limit(30).Find(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}

func FormatVideoInfo(video *Video, existed bool) map[string]interface{} {
	return map[string]interface{}{
		KeyId:       video.Id,
		KeyUserId:   video.UserId,
		KeyTitle:    video.Title,
		KeyPlayUrl:  video.PlayUrl,
		KeyCoverUrl: video.CoverUrl,
		KeyExisted:  existed,
	}
}

func ParseVideoInfo(videoData map[string]string) *Video {
	id, _ := strconv.ParseInt(videoData[KeyId], 10, 64)
	userId, _ := strconv.ParseInt(videoData[KeyUserId], 10, 64)
	return &Video{
		Id:       id,
		UserId:   userId,
		Title:    videoData[KeyTitle],
		PlayUrl:  videoData[KeyPlayUrl],
		CoverUrl: videoData[KeyCoverUrl],
	}
}

func generateExpireTime() time.Duration {
	// 引入随机数减轻缓存雪崩
	return time.Duration(60*60*24+rand.Intn(3600)) * time.Second
}
