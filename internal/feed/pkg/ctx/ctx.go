package ctx

import (
	"toktik/internal/feed/pkg/ctx/feed"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	FeedService *feed.FeedService
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	return &ServiceContext{
		FeedService: feed.NewFeedService(db.Instance),
	}
}
