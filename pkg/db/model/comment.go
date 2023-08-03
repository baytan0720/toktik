package model

import (
	"time"
)

type Comment struct {
	Id        int64     `gorm:"primary_key;auto_increment"`
	VideoId   int64     `gorm:"index,not null"`
	UserId    int64     `gorm:"not null"`
	Content   string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"index,not null"`
}
