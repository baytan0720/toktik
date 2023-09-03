package comment

import (
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type Comment = model.Comment

type CommentService struct {
	dbInstance func() *gorm.DB
}

func NewCommentService(db func() *gorm.DB) *CommentService {
	return &CommentService{
		dbInstance: db,
	}
}

func (c *CommentService) CreateComment(videoId int64, userId int64, content string) (*Comment, error) {
	db := c.dbInstance()
	comment := &Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: content,
	}
	err := db.Create(comment).Error
	return comment, err
}

func (c *CommentService) GetComment(commentId int64) (*Comment, error) {
	db := c.dbInstance()
	comment := &Comment{}
	err := db.Where("id = ?", commentId).First(comment).Error
	return comment, err
}

func (c *CommentService) DeleteComment(commentId int64) error {
	db := c.dbInstance()
	return db.Where("id = ?", commentId).Delete(&Comment{}).Error
}

func (c *CommentService) ListCommentOrderByCreatedAtDesc(videoId int64) ([]*Comment, error) {
	db := c.dbInstance()
	var comments []*Comment
	err := db.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (c *CommentService) CountComment(videoId int64) (int64, error) {
	db := c.dbInstance()
	var count int64
	err := db.Model(&Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}
