package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/internal/comment/kitex_gen/comment"
	commentsvc "toktik/internal/comment/pkg/comment"
	"toktik/internal/comment/pkg/ctx"
	"toktik/internal/user/kitex_gen/user"
	mock_userservice "toktik/pkg/test/mock/user"
	"toktik/pkg/test/testutil"
)

func newMockCommentService(t *testing.T) *commentsvc.CommentService {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&commentsvc.Comment{}))
	return commentsvc.NewCommentService(func() *gorm.DB {
		return db
	})
}

func newMockUserClient(t *testing.T) *mock_userservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_userservice.NewMockClient(ctl)
}

func TestCommentServiceImpl_CreateComment(t *testing.T) {
	t.Run("create comment failed", func(t *testing.T) {
		db := testutil.NewMockDB()
		svcCtx := &ctx.ServiceContext{
			CommentService: commentsvc.NewCommentService(func() *gorm.DB {
				return db
			}),
		}
		svc := NewCommentServiceImpl(svcCtx)
		_, err := svc.CreateComment(context.Background(), &comment.CreateCommentReq{
			UserId:  1,
			VideoId: 1,
			Content: "any content",
		})
		assert.Error(t, err)
	})

	t.Run("get user info failed", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			CommentService: newMockCommentService(t),
			UserClient:     mockUserClient,
		}
		svc := NewCommentServiceImpl(svcCtx)
		resp, err := svc.CreateComment(context.Background(), &comment.CreateCommentReq{
			UserId:  1,
			VideoId: 1,
			Content: "any content",
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.Comment.Content, "any content")
		assert.Equal(t, resp.Comment.CreateDate, time.Now().Format("01-02"))
		assert.Nil(t, resp.Comment.User)
	})

	t.Run("success", func(t *testing.T) {
		mockUserClient := newMockUserClient(t)
		mockUserClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&user.GetUserInfoRes{User: &user.UserInfo{}}, nil).AnyTimes()
		svcCtx := &ctx.ServiceContext{
			CommentService: newMockCommentService(t),
			UserClient:     mockUserClient,
		}
		svc := NewCommentServiceImpl(svcCtx)
		resp, err := svc.CreateComment(context.Background(), &comment.CreateCommentReq{
			UserId:  1,
			VideoId: 1,
			Content: "any content",
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.Comment.Content, "any content")
		assert.Equal(t, resp.Comment.CreateDate, time.Now().Format("01-02"))
		assert.NotNil(t, resp.Comment.User)
	})
}
