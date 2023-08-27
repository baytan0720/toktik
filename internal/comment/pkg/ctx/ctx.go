package ctx

import (
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"

	"toktik/internal/comment/pkg/comment"
	user "toktik/internal/user/kitex_gen/user/userservice"
	"toktik/pkg/config"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	CommentService *comment.CommentService
	UserClient     user.Client
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	r, err := consul.NewConsulResolver(config.Conf.Get("consul").(string))
	if err != nil {
		log.Fatalln(err)
	}

	return &ServiceContext{
		CommentService: comment.NewCommentService(db.Instance),
		UserClient:     user.MustNewClient("user", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}
