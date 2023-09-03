package video

import (
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type Video = model.Video

type VideoService struct {
	dbInstance func() *gorm.DB
}

func NewVideoService(db func() *gorm.DB) *VideoService {
	return &VideoService{
		dbInstance: db,
	}
}

func (s *VideoService) CreateVideo(userId int64, title, playUrl, coverUrl string) error {
	db := s.dbInstance()

	return db.Create(&Video{
		UserId:   userId,
		Title:    title,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
	}).Error
}

func (s *VideoService) ListVideoByUserId(userId int64) ([]*Video, error) {
	db := s.dbInstance()

	videos := make([]*Video, 0)
	if err := db.Where("user_id = ?", userId).Find(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}

func (s *VideoService) GetVideoByIds(videoIds []int64) ([]*Video, error) {
	db := s.dbInstance()

	videos := make([]*Video, 0)
	if err := db.Where("id IN ?", videoIds).Find(&videos).Error; err != nil {
		return nil, err
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
