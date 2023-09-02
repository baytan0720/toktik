package ctx

import (
	favorite "toktik/internal/favorite/kitex_gen/favorite/favoriteservice"
	relation "toktik/internal/relation/kitex_gen/relation/relationservice"
	"toktik/internal/user/pkg/user"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	UserService    *user.UserService
	RelationClient relation.Client
	FavoriteClient favorite.Client
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	return &ServiceContext{
		UserService: user.NewUserService(db.Instance),
	}
}
