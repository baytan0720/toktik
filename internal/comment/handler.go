package main

import (
	"context"

	"toktik/internal/comment/kitex_gen/comment"
	"toktik/internal/comment/pkg/ctx"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct {
	SvcCtx *ctx.ServiceContext
}

func NewCommentServiceImpl(svcCtx *ctx.ServiceContext) *CommentServiceImpl {
	return &CommentServiceImpl{
		SvcCtx: svcCtx,
	}
}

// CreateComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CreateCommentReq) (resp *comment.CreateCommentRes, err error) {
	// TODO: Your code here...
	return
}

// DeleteComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, req *comment.DeleteCommentReq) (resp *comment.DeleteCommentRes, err error) {
	// TODO: Your code here...
	return
}

// ListComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ListComment(ctx context.Context, req *comment.ListCommentReq) (resp *comment.ListCommentRes, err error) {
	// TODO: Your code here...
	return
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, req *comment.GetCommentCountReq) (resp *comment.GetCommentCountRes, err error) {
	// TODO: Your code here...
	return
}
