package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"toktik/internal/gateway/api/user"
	"toktik/internal/gateway/pkg/apiutil"
)

func Register(r *server.Hertz) {
	apiutil.AddRouters(r, user.NewUserAPI())
}
