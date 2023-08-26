package main

import (
	"flag"
	"log"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"

	feed "toktik/internal/feed/kitex_gen/feed/feedservice"
	"toktik/internal/feed/pkg/ctx"
	"toktik/pkg/config"
)

var (
	consulAddr string
	configPath string
)

func init() {
	flag.StringVar(&consulAddr, "consul", "47.115.209.46:8500", "consul address")
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()
}

func main() {
	var conf config.Config
	if configPath != "" {
		conf = config.ReadConfigFromLocal(configPath)
	} else {
		conf = config.ReadConfigFromConsul(consulAddr)
	}
	conf.Set("name", "feed")

	r, err := consul.NewConsulRegister(consulAddr)
	if err != nil {
		log.Fatalf("connect to consul failed: %v", err)
	}

	svr := feed.NewServer(NewFeedServiceImpl(ctx.NewServiceContext()), server.WithRegistry(r), server.WithRegistryInfo(&registry.Info{
		ServiceName: conf.Get("name").(string),
		Tags: map[string]string{
			"release": conf.Get("release").(string),
		},
	}))
	if err := svr.Run(); err != nil {
		log.Println(err)
	}
}
