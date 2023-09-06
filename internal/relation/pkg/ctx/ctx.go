package ctx

import (
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"

	message "toktik/internal/message/kitex_gen/message/messageservice"
	"toktik/internal/relation/pkg/relation"
	user "toktik/internal/user/kitex_gen/user/userservice"
	"toktik/pkg/config"
	"toktik/pkg/db"
)

// ServiceContext contains the components required by the service.
type ServiceContext struct {
	RelationService *relation.RelationService
	UserClient      user.Client
	MessageClient   message.Client
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	db.Init()
	r, err := consul.NewConsulResolver(config.Conf.GetString(config.KEY_CONSUL))
	if err != nil {
		log.Fatalln(err)
	}

	return &ServiceContext{
		RelationService: relation.NewRelationService(db.Instance),
		UserClient:      user.MustNewClient("user", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
		MessageClient:   message.MustNewClient("message", client.WithResolver(r), client.WithRPCTimeout(time.Second*3)),
	}
}
