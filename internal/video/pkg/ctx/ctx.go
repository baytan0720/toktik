package ctx

import (
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/minio/minio-go/v6"

	comment "toktik/internal/comment/kitex_gen/comment/commentservice"
	favorite "toktik/internal/favorite/kitex_gen/favorite/favoriteservice"
	user "toktik/internal/user/kitex_gen/user/userservice"
	"toktik/internal/video/pkg/video"
	"toktik/pkg/config"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	VideoService   *video.VideoService
	MinioClient    *minio.Client
	UserClient     user.Client
	FavoriteClient favorite.Client
	CommentClient  comment.Client
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	minioClient, err := minio.New(
		config.Conf.Get(config.KEY_MINIO_ENDPOINT).(string),
		config.Conf.Get(config.KEY_MINIO_ACCESS_KEY).(string),
		config.Conf.Get(config.KEY_MINIO_SECRET_KEY).(string),
		false,
	)
	if err != nil {
		log.Fatalln("connect to minio failed:", err)
	}
	r, err := consul.NewConsulResolver(config.Conf.Get("consul").(string))
	if err != nil {
		log.Fatalln(err)
	}

	return &ServiceContext{
		VideoService:   video.NewVideoService(db.Instance),
		MinioClient:    minioClient,
		UserClient:     user.MustNewClient("user", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		FavoriteClient: favorite.MustNewClient("favorite", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		CommentClient:  comment.MustNewClient("comment", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}
