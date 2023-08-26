package ctx

import (
	"toktik/internal/video/pkg/video"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	VideoService *video.VideoService
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	return &ServiceContext{
		VideoService: video.NewVideoService(db.Instance),
	}
}
