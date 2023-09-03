package feed

import (
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type Feed = model.Video

type FeedService struct {
	dbInstance func() *gorm.DB
}

func NewFeedService(db func() *gorm.DB) *FeedService {
	return &FeedService{
		dbInstance: db,
	}
}

func (f *FeedService) GetFeed(userId int64, latestTime int64) (feeds []Feed, err error) {
	db := f.dbInstance()

	err = db.Where("user_id = ? AND created_at > ?", userId, latestTime).
		Order("created_at desc").Find(&feeds).Error

	return
}

func (f *FeedService) CheckIsFavorite(userId int64, videoId int64) (isFavorite bool, err error) {
	db := f.dbInstance()
	var favorite bool
	err = db.Table("video").
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Select("is_favorite").
		First(&favorite).Error
	return
}
