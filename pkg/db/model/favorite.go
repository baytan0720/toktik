package model

import "time"

type Favorite struct {
	Id         int64     `gorm:"primary_key;auto_increment"`
	UserId     int64     `gorm:"uniqueIndex:idx_user_id_video_id;not null"`
	VideoId    int64     `gorm:"uniqueIndex:idx_user_id_video_id;index;not null"`
	IsFavorite bool      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}
