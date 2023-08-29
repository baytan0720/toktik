package middleware

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"

	"toktik/internal/gateway/pkg/jwtutil"
)

func TestAuthCheck(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		claims := jwtutil.CreateClaims(1)
		token, err := jwtutil.NewJwtUtil().GenerateToken(claims)
		assert.NoError(t, err)

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("token=" + token)
		AuthCheck(context.Background(), ctx)
		assert.Equal(t, 200, ctx.Response.StatusCode())
		assert.Equal(t, int64(1), ctx.GetInt64(CTX_USER_ID))
	})

	t.Run("invalid token", func(t *testing.T) {
		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("token=" + "invalid_token")
		AuthCheck(context.Background(), ctx)
		assert.Equal(t, 401, ctx.Response.StatusCode())
	})

	t.Run("empty token", func(t *testing.T) {
		ctx := app.NewContext(16)
		AuthCheck(context.Background(), ctx)
		assert.Equal(t, 401, ctx.Response.StatusCode())
	})
}

func TestSoftAuthCheck(t *testing.T) {
	t.Run("with token", func(t *testing.T) {
		claims := jwtutil.CreateClaims(1)
		token, err := jwtutil.NewJwtUtil().GenerateToken(claims)
		assert.NoError(t, err)

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("token=" + token)
		SoftAuthCheck(context.Background(), ctx)
		assert.Equal(t, 200, ctx.Response.StatusCode())
		assert.Equal(t, int64(1), ctx.GetInt64(CTX_USER_ID))
	})

	t.Run("without token", func(t *testing.T) {
		ctx := app.NewContext(16)
		SoftAuthCheck(context.Background(), ctx)
		assert.Equal(t, 200, ctx.Response.StatusCode())
		assert.Equal(t, int64(0), ctx.GetInt64(CTX_USER_ID))
	})
}
