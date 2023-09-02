package ctx

import (
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"time"
	"toktik/internal/favorite/pkg/favorite"
	video "toktik/internal/video/kitex_gen/video/videoservice"
	"toktik/pkg/config"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	FavoriteService *favorite.FavoriteService
	VideoClient     video.Client
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	r, err := consul.NewConsulResolver(config.Conf.Get("consul").(string))
	if err != nil {
		log.Fatalln(err)
	}
	return &ServiceContext{
		FavoriteService: favorite.NewFavoriteService(db.Instance),
		VideoClient:     video.MustNewClient("Video", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}
