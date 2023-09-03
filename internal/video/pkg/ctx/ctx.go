package ctx

import (
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"time"
	"toktik/internal/user/pkg/user"

	"toktik/internal/video/pkg/video"
	"toktik/pkg/config"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	VideoService *video.VideoService
	VideoClient  video.Client
	UserClient   user.Client
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	r, err := consul.NewConsulResolver(config.Conf.Get("consul").(string))
	if err != nil {
		log.Fatalln(err)
	}

	return &ServiceContext{
		VideoService: video.NewVideoService(db.Instance),
		VideoClient:  video.MustNewClient("video", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		UserClient:   user.MustNewClient("user", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}
