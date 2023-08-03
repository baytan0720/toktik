package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"toktik/internal/gateway/api/comment"
	"toktik/internal/gateway/api/favorite"
	"toktik/internal/gateway/api/feed"
	"toktik/internal/gateway/api/message"
	"toktik/internal/gateway/api/publish"
	"toktik/internal/gateway/api/relation"
	"toktik/internal/gateway/api/user"
	"toktik/internal/gateway/pkg/apiutil"
)

func Register(r *server.Hertz) {
	apiutil.AddRouters(r, user.NewUserAPI())
	apiutil.AddRouters(r, comment.NewCommentApi())
	apiutil.AddRouters(r, favorite.NewFavoriteApi())
	apiutil.AddRouters(r, feed.NewFeedApi())
	apiutil.AddRouters(r, message.NewMessageApi())
	apiutil.AddRouters(r, relation.NewRelationApi())
	apiutil.AddRouters(r, publish.NewPublishApi())
}
