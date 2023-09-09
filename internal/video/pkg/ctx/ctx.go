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
	"toktik/pkg/rabbitmq"
	"toktik/pkg/redis"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	VideoService   *video.VideoService
	MinioClient    *minio.Client
	UserClient     user.Client
	FavoriteClient favorite.Client
	CommentClient  comment.Client
	MQ             *rabbitmq.RabbitMQ
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	redis.InitRedisClient()
	minioClient, err := minio.New(
		config.Conf.GetString(config.KEY_MINIO_ENDPOINT),
		config.Conf.GetString(config.KEY_MINIO_ACCESS_KEY),
		config.Conf.GetString(config.KEY_MINIO_SECRET_KEY),
		false,
	)
	if err != nil {
		log.Fatalln("connect to minio failed:", err)
	}
	r, err := consul.NewConsulResolver(config.GetString(config.KEY_CONSUL))
	if err != nil {
		log.Fatalln(err)
	}

	mq, err := rabbitmq.NewProvider(
		config.GetString(config.KEY_RABBITMQ_HOST),
		config.GetString(config.KEY_RABBITMQ_PORT),
		config.GetString(config.KEY_RABBITMQ_USER),
		config.GetString(config.KEY_RABBITMQ_PASSWORD),
		config.GetString(config.KEY_RABBITMQ_QUEUE),
		config.GetString(config.KEY_RABBITMQ_EXCHANGE),
		config.GetString(config.KEY_RABBITMQ_ROUTING_KEY),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return &ServiceContext{
		VideoService: video.NewVideoService(db.Instance, redis.Instance),
		MinioClient:    minioClient,
		UserClient:     user.MustNewClient("user", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		FavoriteClient: favorite.MustNewClient("favorite", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		CommentClient:  comment.MustNewClient("comment", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		MQ:             mq,
	}
}
