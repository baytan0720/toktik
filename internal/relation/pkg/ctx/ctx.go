package ctx

import (
	"toktik/internal/relation/pkg/relation"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	RelationService *relation.RelationService
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	return &ServiceContext{
		RelationService: relation.NewRelationService(db.Instance),
	}
}
