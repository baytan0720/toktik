package ctx

import (
	"toktik/internal/favorite/pkg/favorite"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	FavoriteService *favorite.FavoriteService
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	return &ServiceContext{
		FavoriteService: favorite.NewFavoriteService(db.Instance),
	}
}
