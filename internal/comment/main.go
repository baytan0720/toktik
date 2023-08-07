package main

import (
	"flag"
	"log"

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
	flag.Parse()
}

func main() {
	var conf config.Config
	if configPath != "" {
		conf = config.ReadConfigFromLocal(configPath)
	} else {
		conf = config.ReadConfigFromConsul(consulAddr)
	}
	conf.Set("name", "comment")

	svr := comment.NewServer(NewCommentServiceImpl(ctx.NewServiceContext()))
	if err := svr.Run(); err != nil {
		log.Println(err)
	}
}
