package apiutil

import (
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type APIRoutes interface {
	Routes() []Route
}

type Route struct {
	Method  string
	Path    string
	Hooks   []app.HandlerFunc
	Handler app.HandlerFunc
}

func AddRouters(r *server.Hertz, api APIRoutes) {
	for _, route := range api.Routes() {
		handlers := append(route.Hooks, route.Handler)

		if route.Method != "" {
			r.Handle(route.Method, route.Path, handlers...)
		} else {
			r.Any(route.Path, handlers...)
		}
	}
}

type H map[string]any

func HandleError(ctx *app.RequestContext, err, errMsg error) bool {
	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusOK, H{
			"status_code": StatusFailed,
			"status_msg":  errMsg.Error(),
		})
		return true
	}
	return false
}

func HandleRpcError(ctx *app.RequestContext, status int32, errMsg string) bool {
	if status != 0 {
		ctx.JSON(http.StatusOK, H{
			"status_code": StatusFailed,
			"status_msg":  errMsg,
		})
		return true
	}
	return false
}
