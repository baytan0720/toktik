package middleware

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/gateway/pkg/jwtutil"
)

const CTX_USER_ID = "user_id"

func AuthCheck(c context.Context, ctx *app.RequestContext) {
	token := ctx.Query("token")
	claims, err := jwtutil.ParseToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
			"status_code": apiutil.StatusFailed,
			"status_msg":  err.Error(),
		})
		return
	}
	ctx.Set(CTX_USER_ID, claims.UserId)
	ctx.Next(c)
}

func SoftAuthCheck(c context.Context, ctx *app.RequestContext) {
	token := ctx.Query("token")
	claims, err := jwtutil.ParseToken(token)
	if err != nil {
		ctx.Set(CTX_USER_ID, 0)
	} else {
		ctx.Set(CTX_USER_ID, claims.UserId)
	}
	ctx.Next(c)
}
