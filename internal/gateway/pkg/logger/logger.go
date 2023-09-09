package logger

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/sirupsen/logrus"

	"toktik/pkg/config"
)

func InitLogger() {
	logger := hertzlogrus.NewLogger(hertzlogrus.WithLogger(logrus.New()))
	hlog.SetLogger(logger)
	setLevel()
	setOutput()
}

func setLevel() {
	level := config.GetString(config.KEY_LOGGER_LEVEL)
	switch level {
	case "debug":
		hlog.SetLevel(hlog.LevelDebug)
	case "info":
		hlog.SetLevel(hlog.LevelInfo)
	case "warn":
		hlog.SetLevel(hlog.LevelWarn)
	case "error":
		hlog.SetLevel(hlog.LevelError)
	default:
		hlog.SetLevel(hlog.LevelInfo)
	}
}

func setOutput() {
	output := config.GetString(config.KEY_LOGGER_OUTPUT)
	if output != "" {
		f, err := os.OpenFile(path.Join(output, time.Now().Format("2006-01-02")+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			hlog.Fatal("failed to open log file:", err)
		}
		hlog.SetOutput(f)
	}
}

func LoggerHandler() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		hlog.Infof("%d\t%s\t%s", ctx.Response.StatusCode(), ctx.Request.Method(), ctx.Request.RequestURI())
	}
}
