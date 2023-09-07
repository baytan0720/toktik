package main

import (
	"flag"
	"log"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"

	comment "toktik/internal/comment/kitex_gen/comment/commentservice"
	"toktik/internal/comment/pkg/ctx"
	"toktik/pkg/config"
)

var (
	consulAddr string
	configPath string
)

func init() {
	flag.StringVar(&consulAddr, "consul", "47.115.209.46:8500", "consul address")
	flag.StringVar(&configPath, "config", "", "config path")
}

func main() {
	flag.Parse()
	var conf config.Config
	if configPath != "" {
		conf = config.ReadConfigFromLocal(configPath)
	} else {
		conf = config.ReadConfigFromConsul(consulAddr)
	}
	conf.Set(config.KEY_SERVICE_NAME, "comment")
	conf.Set(config.KEY_CONSUL, consulAddr)

	r, err := consul.NewConsulRegister(consulAddr)
	if err != nil {
		log.Fatalf("connect to consul failed: %v", err)
	}

	svr := comment.NewServer(NewCommentServiceImpl(ctx.NewServiceContext()), server.WithRegistry(r), server.WithRegistryInfo(&registry.Info{
		ServiceName: conf.GetString(config.KEY_SERVICE_NAME),
		Tags: map[string]string{
			"release": conf.GetString(config.KEY_RELEASE),
		},
	}))
	if err := svr.Run(); err != nil {
		log.Println(err)
	}
}
