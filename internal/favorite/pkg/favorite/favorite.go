package favorite

import (
	"gorm.io/gorm"
)

type Favorite = any

type FavoriteService struct {
	dbInstance func() *gorm.DB
}

func NewFavoriteService(db func() *gorm.DB) *FavoriteService {
	return &FavoriteService{
		dbInstance: db,
	}
}
