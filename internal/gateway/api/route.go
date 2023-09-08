package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	consul "github.com/kitex-contrib/registry-consul"

	"toktik/internal/gateway/api/comment"
	"toktik/internal/gateway/api/favorite"
	"toktik/internal/gateway/api/feed"
	"toktik/internal/gateway/api/message"
	"toktik/internal/gateway/api/publish"
	"toktik/internal/gateway/api/relation"
	"toktik/internal/gateway/api/user"
	"toktik/internal/gateway/pkg/apiutil"
	"toktik/internal/gateway/pkg/logger"
	"toktik/pkg/config"
)

func Register(r *server.Hertz) {
	resolver, err := consul.NewConsulResolver(config.GetString(config.KEY_CONSUL))
	if err != nil {
		hlog.Fatal(err)
	}

	apiutil.AddRouters(r, user.NewUserAPI(resolver))
	apiutil.AddRouters(r, comment.NewCommentApi(resolver))
	apiutil.AddRouters(r, favorite.NewFavoriteApi(resolver))
	apiutil.AddRouters(r, feed.NewFeedApi(resolver))
	apiutil.AddRouters(r, message.NewMessageApi(resolver))
	apiutil.AddRouters(r, relation.NewRelationApi(resolver))
	apiutil.AddRouters(r, publish.NewPublishApi(resolver))

	r.Use(logger.LoggerHandler())
}
