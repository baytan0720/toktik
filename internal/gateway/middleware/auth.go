package middleware

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"toktik/internal/gateway/pkg/jwtutil"
)

const CTX_USER_ID = "user_id"

func AuthCheck(c context.Context, ctx *app.RequestContext) {
	token := ctx.Query("token")
	jwtUtil := jwtutil.NewJwtUtil()
	claims, err := jwtUtil.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status_msg": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.Set(CTX_USER_ID, claims.UserId)
	ctx.Next(c)
}

func SoftAuthCheck(c context.Context, ctx *app.RequestContext) {
	token := ctx.Query("token")
	jwtUtil := jwtutil.NewJwtUtil()
	claims, err := jwtUtil.ParseToken(token)
	if err != nil {
		ctx.Set(CTX_USER_ID, 0)
	} else {
		ctx.Set(CTX_USER_ID, claims.UserId)
	}
	ctx.Next(c)
}
