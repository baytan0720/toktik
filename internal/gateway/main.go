package main

import (
	"flag"

	"github.com/cloudwego/hertz/pkg/app/server"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/cors"

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
	conf.Set(config.KEY_SERVICE_NAME, "gateway")
	conf.Set(config.KEY_CONSUL, consulAddr)

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulAddr

	router := server.Default(
		server.WithHostPorts(conf.GetString(config.KEY_LISTEN_ON)),
		server.WithKeepAlive(true),
	)

	router.Use(cors.Default())

	api.Register(router)

	router.Spin()
}
