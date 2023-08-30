package model

import (
	"time"
)

type Video struct {
	Id        int64     `gorm:"primary_key;auto_increment"`
	UserId    int64     `gorm:"index;not null"`
	Title     string    `gorm:"type:varchar(255);not null"`
	PlayUrl   string    `gorm:"type:varchar(255);not null"`
	CoverUrl  string    `gorm:"type:varchar(255); not null"`
	CreatedAt time.Time `gorm:"index;not null"`
}
