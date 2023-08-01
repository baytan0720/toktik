package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"toktik/internal/gateway/api"
)

func main() {
	r := server.Default()

	api.Register(r)

	r.Spin()
}
