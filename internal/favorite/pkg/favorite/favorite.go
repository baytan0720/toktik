package favorite

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"toktik/pkg/db/model"
)

type Favorite = model.Favorite

type FavoriteService struct {
	dbInstance func() *gorm.DB
}

func NewFavoriteService(db func() *gorm.DB) *FavoriteService {
	return &FavoriteService{
		dbInstance: db,
	}
}

func (f *FavoriteService) Favorite(videoId int64, userId int64) (err error) {
	db := f.dbInstance()

	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "video_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"is_favorite": true, "updated_at": time.Now()}),
	}).Create(&Favorite{
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: true,
	}).Error

	return
}

func (f *FavoriteService) UnFavorite(videoId int64, userId int64) error {
	db := f.dbInstance()

	var favorite Favorite
	err := db.Where("user_id = ? AND video_id = ? AND is_favorite = ?", videoId, userId, true).First(&favorite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	favorite.IsFavorite = false
	err = db.Save(&favorite).Error

	return err
}

func (f *FavoriteService) ListFavoriteByUserId(userId int64) ([]int64, error) {
	db := f.dbInstance()

	var videoIdList []int64
	err := db.Model(&Favorite{}).Select("video_id").Where("user_id = ? AND is_favorite = ?", userId, true).Order("updated_at desc").Find(&videoIdList).Error
	if err != nil {
		return nil, err
	}

	return videoIdList, nil
}

func (f *FavoriteService) CountVideoFavorite(videoId int64) (int64, error) {
	db := f.dbInstance()

	var count int64
	err := db.Model(&Favorite{}).Where("video_id = ? AND is_favorite = ?", videoId, true).Count(&count).Error

	return count, err
}

func (f *FavoriteService) CountUserFavorite(userId int64) (int64, error) {
	db := f.dbInstance()

	var count int64
	err := db.Model(&Favorite{}).Where("user_id = ? AND is_favorite = ?", userId, true).Count(&count).Error

	return count, err
}
