package main

import (
	"context"
	"errors"
	"strings"

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

var (
	errCommentFailed       = "评论失败"
	errCommentNotExists    = "评论不存在"
	errCommentUserNotMatch = "没有权限"
	errDeleteCommentFailed = "删除评论失败"
	errListCommentFailed   = "获取评论列表失败"
)

// CreateComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CreateCommentReq) (resp *comment.CreateCommentRes, _ error) {
	resp = &comment.CreateCommentRes{}
	commentInfo, err := s.svcCtx.CommentService.CreateComment(req.VideoId, req.UserId, req.Content)
	if err != nil {
		resp.Status = comment.Status_ERROR
		resp.ErrMsg = errCommentFailed
		return
	}

	resp.Comment = &comment.CommentInfo{
		Id:         commentInfo.Id,
		Content:    commentInfo.Content,
		CreateDate: commentInfo.CreatedAt.Format("01-02"),
	}

	res, err := s.svcCtx.UserClient.GetUserInfo(ctx, &user.GetUserInfoReq{
		ToUserId: req.UserId,
	})
	if err != nil {
		return
	}

	if res.Status == user.Status_OK {
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

		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.ErrMsg = errCommentNotExists
			return
		}

		if strings.Contains(err.Error(), "not match") {
			resp.ErrMsg = errCommentUserNotMatch
			return
		}

		resp.ErrMsg = errDeleteCommentFailed
	}
	return
}

// ListComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ListComment(ctx context.Context, req *comment.ListCommentReq) (resp *comment.ListCommentRes, _ error) {
	resp = &comment.ListCommentRes{}

	commentInfos, err := s.svcCtx.CommentService.ListCommentOrderByCreatedAtDesc(req.VideoId)
	if err != nil {
		resp.Status = comment.Status_ERROR
		resp.ErrMsg = errListCommentFailed
		return
	}

	toUserIdList := make([]int64, 0, len(commentInfos))
	for _, c := range commentInfos {
		commentInfo := &comment.CommentInfo{
			Id:         c.Id,
			Content:    c.Content,
			CreateDate: c.CreatedAt.Format("01-02"),
		}
		resp.CommentList = append(resp.CommentList, commentInfo)
		toUserIdList = append(toUserIdList, c.UserId)
	}

	res, err := s.svcCtx.UserClient.GetUserInfos(ctx, &user.GetUserInfosReq{
		UserId:    req.UserId,
		ToUserIds: toUserIdList,
	})
	if err != nil {
		return
	}

	id2User := make(map[int64]*comment.UserInfo)
	if res.Status == user.Status_OK {
		for _, u := range res.Users {
			id2User[u.Id] = convert2CommentUserInfo(u)
		}
	}

	for _, c := range resp.CommentList {
		if u, ok := id2User[c.Id]; ok {
			c.User = u
		}
	}

	return
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, req *comment.GetCommentCountReq) (resp *comment.GetCommentCountRes, _ error) {
	resp = &comment.GetCommentCountRes{}

	videoIdList := req.VideoIdList

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
