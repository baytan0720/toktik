package model

type video struct {
	Id            int64 `gorm:"primary_key;AUTO_INCREMENT"`
	Author        User
	PlayUrl       string `gorm:"type:varchar(128);not null"`
	CoverUrl      string `gorm:"type:varchar(128);not null"`
	FavoriteCount int64  `gorm:"default:0"`
	CommentCount  int64  `gorm:"default:0"`
	IsFavorite    bool   `gorm:"default:null"`
	Title         string `gorm:"type:varchar(128);not null"`
}
