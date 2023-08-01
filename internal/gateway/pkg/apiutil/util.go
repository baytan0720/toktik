package apiutil

import (
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
