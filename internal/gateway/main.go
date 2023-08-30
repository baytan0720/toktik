package main

import (
	"flag"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"

	"toktik/internal/gateway/api"
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
	conf.Set("name", "gateway")
	conf.Set("consul", consulAddr)

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulAddr
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("connect to consul failed: %v", err)
	}
	r := consul.NewConsulRegister(consulClient)

	router := server.Default(server.WithRegistry(r, &registry.Info{
		ServiceName: conf.Get("name").(string),
		Tags: map[string]string{
			"release": conf.Get("release").(string),
		},
	}))

	api.Register(router)

	router.Spin()
}
