package ctx

import (
	"toktik/internal/user/pkg/user"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	UserService *user.UserService
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	return &ServiceContext{
		UserService: user.NewUserService(db.Instance),
	}
}
