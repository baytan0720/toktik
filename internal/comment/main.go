package main

import (
	"log"

	comment "toktik/internal/comment/kitex_gen/comment/commentservice"
	"toktik/internal/comment/pkg/ctx"
)

func main() {
	svr := comment.NewServer(NewCommentServiceImpl(ctx.NewServiceContext()))
	if err := svr.Run(); err != nil {
		log.Println(err)
	}
}
