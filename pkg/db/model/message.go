package model

import (
	"time"
)

type Message struct {
	Id             int64     `gorm:"primary_key;auto_increment"`
	UserId         int64     `gorm:"index;not null"`
	ToUserId       int64     `gorm:"index;not null"`
	Content        string    `gorm:"type:varchar(255);not null"`
	LastMessage    string    `gorm:"type:varchar(255);not null"`
	MessageType    int64     `gorm:"default:0;not null"`
	CreatedAt      time.Time `gorm:"index;not null"`
	PreMessageTime int64     `gorm:"not null"` // 上次最新消息的时间
}
