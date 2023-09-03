package video

import (
	"gorm.io/gorm"
	"time"
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
func (v *VideoService) ListVideo(userId int64) ([]Video, error) {
	db := v.dbInstance()
	var videos []Video
	err := db.Where("user_id = ?", userId).Find(&videos).Error
	return videos, err
}

func (v *VideoService) GetVideoCount(userId int64) (int64, error) {
	db := v.dbInstance()
	var count int64
	err := db.Where("user_id=?", userId).Count(&count).Error
	return count, err
}

func (v *VideoService) GetVideo(videoId int64) (Video, error) {
	db := v.dbInstance()
	var video Video
	err := db.Where("id=?", videoId).First(&video).Error
	return video, err
}

func (v *VideoService) PublishVideo(videoName string, imageName string, userId int64, title string) error {
	db := v.dbInstance()
	var video Video
	video.CreatedAt = time.Now()
	video.UserId = userId
	video.Title = title
	video.PlayUrl = videoName + ".mp4"
	video.CoverUrl = imageName + ".jpg"
	err := db.Save(&video).Error
	return err
}
