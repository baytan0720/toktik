package api

import (
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
	consul "github.com/kitex-contrib/registry-consul"

	"toktik/internal/gateway/api/comment"
	"toktik/internal/gateway/api/favorite"
	"toktik/internal/gateway/api/feed"
	"toktik/internal/gateway/api/message"
	"toktik/internal/gateway/api/publish"
	"toktik/internal/gateway/api/relation"
	"toktik/internal/gateway/api/user"
	"toktik/internal/gateway/pkg/apiutil"
	"toktik/pkg/config"
)

func Register(r *server.Hertz) {
	resolver, err := consul.NewConsulResolver(config.Conf.Get("consul").(string))
	if err != nil {
		log.Fatalln(err)
	}

	apiutil.AddRouters(r, user.NewUserAPI(resolver))
	apiutil.AddRouters(r, comment.NewCommentApi())
	apiutil.AddRouters(r, favorite.NewFavoriteApi())
	apiutil.AddRouters(r, feed.NewFeedApi(resolver))
	apiutil.AddRouters(r, message.NewMessageApi())
	apiutil.AddRouters(r, relation.NewRelationApi())
	apiutil.AddRouters(r, publish.NewPublishApi(resolver))
}
