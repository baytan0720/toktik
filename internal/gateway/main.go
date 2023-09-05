package main

import (
	"flag"

	"github.com/cloudwego/hertz/pkg/app/server"
	consulapi "github.com/hashicorp/consul/api"

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

	router := server.Default(
		server.WithHostPorts("0.0.0.0:8888"),
	)

	api.Register(router)

	router.Spin()
}
