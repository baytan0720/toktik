package ctx

import (
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"

	favorite "toktik/internal/favorite/kitex_gen/favorite/favoriteservice"
	relation "toktik/internal/relation/kitex_gen/relation/relationservice"
	"toktik/internal/user/pkg/user"
	video "toktik/internal/video/kitex_gen/video/videoservice"
	"toktik/pkg/config"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	UserService    *user.UserService
	FavoriteClient favorite.Client
	VideoClient    video.Client
	RelationClient relation.Client
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	r, err := consul.NewConsulResolver(config.Conf.Get("consul").(string))
	if err != nil {
		log.Fatalln(err)
	}

	return &ServiceContext{
		UserService:    user.NewUserService(db.Instance),
		FavoriteClient: favorite.MustNewClient("favorite", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		VideoClient:    video.MustNewClient("video", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		RelationClient: relation.MustNewClient("relation", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}
