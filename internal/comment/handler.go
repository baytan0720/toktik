package main

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"toktik/internal/comment/kitex_gen/comment"
	"toktik/internal/comment/pkg/ctx"
	"toktik/internal/user/kitex_gen/user"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewCommentServiceImpl(svcCtx *ctx.ServiceContext) *CommentServiceImpl {
	return &CommentServiceImpl{
		svcCtx: svcCtx,
	}
}

// CreateComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CreateCommentReq) (resp *comment.CreateCommentRes, err error) {
	resp = &comment.CreateCommentRes{}
	commentInfo, err := s.svcCtx.CommentService.CreateComment(req.VideoId, req.UserId, req.Content)
	if err != nil {
		// internal error
		return nil, err
	}

	resp.Comment = &comment.CommentInfo{
		Id:         commentInfo.Id,
		Content:    commentInfo.Content,
		CreateDate: commentInfo.CreatedAt.Format("01-02"),
	}

	if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
		ToUserId: req.UserId,
	}); err == nil {
		resp.Comment.User = convert2CommentUserInfo(res.User)
	} else {
		log.Println("get user info failed:", err)
	}

	return resp, nil
}

// DeleteComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, req *comment.DeleteCommentReq) (resp *comment.DeleteCommentRes, err error) {
	resp = &comment.DeleteCommentRes{}
	err = s.svcCtx.CommentService.DeleteComment(req.UserId, req.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, errors.New("user id not match")) {
			resp.Status = comment.Status_ERROR
			resp.ErrMsg = err.Error()
		} else {
			return nil, err
		}
	}
	return resp, nil
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

func convert2CommentUserInfo(user *user.UserInfo) *comment.UserInfo {
	return &comment.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        false,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
