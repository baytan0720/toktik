package comment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"toktik/internal/comment/kitex_gen/comment"
	"toktik/internal/gateway/pkg/apiutil"
	mock_commentservice "toktik/pkg/test/mock/comment"
)

func newMockCommentClient(t *testing.T) *mock_commentservice.MockClient {
	ctl := gomock.NewController(t)
	return mock_commentservice.NewMockClient(ctl)
}

func TestUserAPI_Action(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockCommentClient(t)
		api := &CommentApi{
			commentClient: mockUserClient,
		}

		mockUserClient.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10&action_type=1&comment_text='that's awesome!'")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("comment failed", func(t *testing.T) {
		mockUserClient := newMockCommentClient(t)
		api := &CommentApi{
			commentClient: mockUserClient,
		}

		mockUserClient.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(&comment.CreateCommentRes{
			Status: comment.Status_ERROR,
			ErrMsg: "fail to comment",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10&action_type=1&comment_text='that's awesome!'")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "fail to comment", payload.StatusMsg)
	})

	t.Run("comment success", func(t *testing.T) {
		mockUserClient := newMockCommentClient(t)
		api := &CommentApi{
			commentClient: mockUserClient,
		}

		info := &comment.CommentInfo{
			Content:    "that's awesome!",
			CreateDate: "string",
			Id:         10,
			User:       nil,
		}
		mockUserClient.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(&comment.CreateCommentRes{
			Status:  comment.Status_OK,
			Comment: info,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10&action_type=1&comment_text='that's awesome!'")
		api.Action(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ActionResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, info, payload.Comment)
	})
}

func TestUserAPI_List(t *testing.T) {
	t.Run("rpc error", func(t *testing.T) {
		mockUserClient := newMockCommentClient(t)
		api := &CommentApi{
			commentClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListComment(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc error")).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())
		payload := ListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "rpc error", payload.StatusMsg)
	})

	t.Run("list failed", func(t *testing.T) {
		mockUserClient := newMockCommentClient(t)
		api := &CommentApi{
			commentClient: mockUserClient,
		}

		mockUserClient.EXPECT().ListComment(gomock.Any(), gomock.Any()).Return(&comment.ListCommentRes{
			Status: comment.Status_ERROR,
			ErrMsg: "fail to list",
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response.StatusCode())
		payload := ListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusFailed, payload.StatusCode)
		assert.Equal(t, "fail to list", payload.StatusMsg)
	})

	t.Run("list success", func(t *testing.T) {
		mockUserClient := newMockCommentClient(t)
		api := &CommentApi{
			commentClient: mockUserClient,
		}

		infos := make([]*comment.CommentInfo, 10)
		mockUserClient.EXPECT().ListComment(gomock.Any(), gomock.Any()).Return(&comment.ListCommentRes{
			Status:      comment.Status_OK,
			CommentList: infos,
		}, nil).AnyTimes()

		ctx := app.NewContext(16)
		ctx.Request.SetQueryString("video_id=10")
		api.List(context.Background(), ctx)

		assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())
		payload := ListResp{}
		assert.NoError(t, json.Unmarshal(ctx.Response.Body(), &payload))
		assert.Equal(t, apiutil.StatusOK, payload.StatusCode)
		assert.Equal(t, infos, payload.CommentList)
	})
}
