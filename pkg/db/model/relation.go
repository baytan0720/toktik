package model

import (
	"time"
)

type FollowRelation struct {
	Id        int64     `gorm:"primary_key;auto_increment"`
	UserId    int64     `gorm:"uniqueIndex:idx_user_id_to_user_id;not null"`
	ToUserId  int64     `gorm:"uniqueIndex:idx_user_id_to_user_id;index;not null"`
	IsFollow  bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
