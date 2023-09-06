package main

import (
	"context"
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
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CreateCommentReq) (resp *comment.CreateCommentRes, _ error) {
	resp = &comment.CreateCommentRes{}
	commentInfo, err := s.svcCtx.CommentService.CreateComment(req.VideoId, req.UserId, req.Content)
	if err != nil {
		// internal error
		resp.Status = comment.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	resp.Comment = &comment.CommentInfo{
		Id:         commentInfo.Id,
		Content:    commentInfo.Content,
		CreateDate: commentInfo.CreatedAt.Format("01-02"),
	}

	if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
		ToUserId: req.UserId,
	}); err != nil {
		resp.Status = comment.Status_ERROR
		resp.ErrMsg = err.Error()
	} else if res.Status == user.Status_OK {
		resp.Comment.User = convert2CommentUserInfo(res.User)
	}

	return
}

// DeleteComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, req *comment.DeleteCommentReq) (resp *comment.DeleteCommentRes, _ error) {
	resp = &comment.DeleteCommentRes{}
	err := s.svcCtx.CommentService.DeleteComment(req.UserId, req.CommentId)
	if err != nil {
		resp.Status = comment.Status_ERROR
		resp.ErrMsg = err.Error()
	}
	return
}

// ListComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ListComment(ctx context.Context, req *comment.ListCommentReq) (resp *comment.ListCommentRes, _ error) {
	resp = &comment.ListCommentRes{}

	commentInfos, err := s.svcCtx.CommentService.ListCommentOrderByCreatedAtDesc(req.VideoId)
	if err != nil {
		resp.Status = comment.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	for _, c := range commentInfos {
		commentInfo := &comment.CommentInfo{
			Id:         c.Id,
			Content:    c.Content,
			CreateDate: c.CreatedAt.Format("01-02"),
		}
		if res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
			ToUserId: c.UserId,
		}); err == nil {
			commentInfo.User = convert2CommentUserInfo(res.User)
		} else {
			resp.Status = comment.Status_ERROR
			resp.ErrMsg = err.Error()
			return
		}
		resp.CommentList = append(resp.CommentList, commentInfo)
	}
	return
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, req *comment.GetCommentCountReq) (resp *comment.GetCommentCountRes, _ error) {
	resp = &comment.GetCommentCountRes{}
	videoIdList := req.VideoIdList
	if len(videoIdList) == 0 {
		return
	}
	for _, videoId := range videoIdList {
		count, err := s.svcCtx.CommentService.CountComment(videoId)
		if err != nil {
			resp.Status = comment.Status_ERROR
			resp.ErrMsg = err.Error()
			return
		}
		resp.CommentCountList = append(resp.CommentCountList, &comment.CommentCountInfo{
			VideoId: videoId,
			Count:   count,
		})
	}

	return
}

func convert2CommentUserInfo(user *user.UserInfo) *comment.UserInfo {
	return &comment.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
