package operator

import (
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type CommentOperator struct {
	dbInstance func() *gorm.DB
}

func NewCommentOperator(db func() *gorm.DB) *CommentOperator {
	return &CommentOperator{
		dbInstance: db,
	}
}

func (o *CommentOperator) CreateComment(videoId int64, userId int64, content string) (*model.Comment, error) {
	db := o.dbInstance()
	comment := &model.Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: content,
	}
	err := db.Create(comment).Error
	return comment, err
}

func (o *CommentOperator) GetComment(commentId int64) (*model.Comment, error) {
	db := o.dbInstance()
	comment := &model.Comment{}
	err := db.Where("id = ?", commentId).First(comment).Error
	return comment, err
}

func (o *CommentOperator) DeleteComment(commentId int64) error {
	db := o.dbInstance()
	return db.Where("id = ?", commentId).Delete(&model.Comment{}).Error
}

func (o *CommentOperator) ListCommentOrderByCreatedAtDesc(videoId int64) ([]*model.Comment, error) {
	db := o.dbInstance()
	comments := []*model.Comment{}
	err := db.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (o *CommentOperator) CountComment(videoId int64) (int64, error) {
	db := o.dbInstance()
	var count int64
	err := db.Model(&model.Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}
