package video

import (
	"gorm.io/gorm"
)

type Video = any

type VideoService struct {
	dbInstance func() *gorm.DB
}

func NewVideoService(db func() *gorm.DB) *VideoService {
	return &VideoService{
		dbInstance: db,
	}
}
