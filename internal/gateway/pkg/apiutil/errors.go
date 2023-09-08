package apiutil

import (
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	ErrInvalidParams = errors.New("参数错误")
	ErrInternalError = errors.New("服务器错误")
)

type H map[string]any

func HandleError(ctx *app.RequestContext, err, errMsg error) bool {
	if err != nil {
		hlog.Warn(err)
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
