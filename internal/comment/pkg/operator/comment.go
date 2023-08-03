package operator

import (
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type CommentOperator struct {
	db *gorm.DB
}

func NewCommentOperator(db *gorm.DB) *CommentOperator {
	return &CommentOperator{
		db: db,
	}
}

func (m *CommentOperator) CreateComment(videoId int64, userId int64, content string) (*model.Comment, error) {
	comment := &model.Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: content,
	}
	err := m.db.Create(comment).Error
	return comment, err
}

func (m *CommentOperator) GetComment(commentId int64) (*model.Comment, error) {
	comment := &model.Comment{}
	err := m.db.Where("id = ?", commentId).First(comment).Error
	return comment, err
}

func (m *CommentOperator) DeleteComment(commentId int64) error {
	return m.db.Where("id = ?", commentId).Delete(&model.Comment{}).Error
}

func (m *CommentOperator) ListCommentOrderByCreatedAtDesc(videoId int64) ([]*model.Comment, error) {
	comments := []*model.Comment{}
	err := m.db.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (m *CommentOperator) CountComment(videoId int64) (int64, error) {
	var count int64
	err := m.db.Model(&model.Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}
