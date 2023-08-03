package ctx

import (
	"toktik/internal/comment/pkg/operator"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	Operator *operator.CommentOperator
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	// init database
	db.Init("")

	return &ServiceContext{
		Operator: operator.NewCommentOperator(db.Instance()),
	}
}
