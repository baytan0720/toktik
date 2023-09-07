package apiutil

import "errors"

var (
	ErrInvalidParams = errors.New("参数错误")
	ErrInternalError = errors.New("服务器错误")
)
