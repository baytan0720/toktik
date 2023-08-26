package ctx

import (
	"toktik/internal/comment/pkg/comment"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	CommentService *comment.CommentService
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	return &ServiceContext{
		CommentService: comment.NewCommentService(db.Instance),
	}
}
