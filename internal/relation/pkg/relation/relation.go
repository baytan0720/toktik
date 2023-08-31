package relation

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"toktik/pkg/db/model"
)

type Relation = model.Relation

type RelationService struct {
	dbInstance func() *gorm.DB
}

func NewRelationService(db func() *gorm.DB) *RelationService {
	return &RelationService{
		dbInstance: db,
	}
}

func (r *RelationService) GetFollowRelations(userId int64, toUserIdList []int64) (relations []Relation, err error) {
	db := r.dbInstance()

	err = db.Where("to_user_id in ?", toUserIdList).Where(&Relation{UserId: userId}).Find(&relations).Error

	return
}

func (r *RelationService) GetFollowCount(userId int64) (count int64, err error) {
	db := r.dbInstance()

	err = db.Model(&Relation{}).Where(&Relation{UserId: userId, IsFollow: true}).Count(&count).Error

	return
}

func (r *RelationService) GetFollowerCount(toUserId int64) (count int64, err error) {
	db := r.dbInstance()

	err = db.Model(&Relation{}).Where(&Relation{ToUserId: toUserId, IsFollow: true}).Count(&count).Error

	return
}

func (r *RelationService) Follow(userId int64, toUserId int64) (err error) {
	db := r.dbInstance()

	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "to_user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"is_follow": true, "updated_at": time.Now()}),
	}).Create(&Relation{
		UserId:   userId,
		ToUserId: toUserId,
		IsFollow: true,
	}).Error

	return
}

func (r *RelationService) Unfollow(userId int64, toUserId int64) (err error) {
	db := r.dbInstance()
	var relation Relation

	err = db.Where(&Relation{UserId: userId, ToUserId: toUserId}).First(&relation).Error
	if err != nil {
		return
	}
	// 记录存在则更新
	relation.IsFollow = false

	err = db.Save(&relation).Error

	return
}

func (r *RelationService) IsFollow(userId int64, toUserId int64) (isFollow bool, err error) {
	db := r.dbInstance()

	var relation Relation

	err = db.Where(&Relation{UserId: userId, ToUserId: toUserId}).First(&relation).Error

	if err == nil {
		isFollow = relation.IsFollow
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return
}

func (r *RelationService) ListFollow(userId int64) (userIdList []int64, err error) {
	db := r.dbInstance()

	err = db.Model(&Relation{}).Select("to_user_id").Where(&Relation{UserId: userId, IsFollow: true}).Order("updated_at desc").Find(&userIdList).Error

	return
}

func (r *RelationService) ListFollower(userId int64) (userIdList []int64, err error) {
	db := r.dbInstance()

	err = db.Model(&Relation{}).Select("user_id").Where(&Relation{ToUserId: userId, IsFollow: true}).Order("updated_at desc").Find(&userIdList).Error

	return
}
