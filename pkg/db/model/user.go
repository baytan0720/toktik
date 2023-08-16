package model

import "time"

type User struct {
	Id              int64     `gorm:"primary_key;AUTO_INCREMENT"`
	Username        string    `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password        string    `gorm:"type:varchar(32);not null"`
	Avatar          string    `gorm:"type:varchar(128)"`
	BackgroundImage string    `gorm:"type:varchar(128)"`
	Signature       string    `gorm:"type:varchar(128)"`
	CreatedAt       time.Time `gorm:"not null"`
}
