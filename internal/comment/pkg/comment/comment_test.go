package comment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/pkg/db/model"
	"toktik/pkg/test/testutil"
)

func newMockDB(t *testing.T) *gorm.DB {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&model.Comment{}))
	return db
}

func TestCommentOperator_CreateComment(t *testing.T) {
	db := newMockDB(t)

	c := NewCommentService(func() *gorm.DB {
		return db
	})
	comment, err := c.CreateComment(1, 1, "test")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), comment.VideoId)
	assert.Equal(t, int64(1), comment.UserId)
	assert.Equal(t, "test", comment.Content)
}

func TestCommentOperator_GetComment(t *testing.T) {
	db := newMockDB(t)

	testCommentCase := &model.Comment{
		VideoId: 10,
		UserId:  10,
		Content: "test comment",
	}
	db.Create(testCommentCase)

	c := NewCommentService(func() *gorm.DB {
		return db
	})
	comment, err := c.GetComment(testCommentCase.Id)
	assert.NoError(t, err)
	assert.Equal(t, testCommentCase.VideoId, comment.VideoId)
	assert.Equal(t, testCommentCase.UserId, comment.UserId)
	assert.Equal(t, testCommentCase.Content, comment.Content)
}

func TestCommentOperator_DeleteComment(t *testing.T) {
	db := newMockDB(t)

	t.Run("success", func(t *testing.T) {
		testCommentCase := &model.Comment{
			VideoId: 10,
			UserId:  10,
			Content: "test comment",
		}
		db.Create(testCommentCase)

		c := NewCommentService(func() *gorm.DB {
			return db
		})
		err := c.DeleteComment(10, testCommentCase.Id)
		assert.NoError(t, err)
		err = db.Where("id = ?", testCommentCase.Id).First(&model.Comment{}).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		testCommentCase := &model.Comment{
			VideoId: 10,
			UserId:  10,
			Content: "test comment",
		}
		db.Create(testCommentCase)

		c := NewCommentService(func() *gorm.DB {
			return db
		})
		err := c.DeleteComment(11, testCommentCase.Id)
		assert.Error(t, err)
	})
}

func TestCommentOperator_ListComment(t *testing.T) {
	db := newMockDB(t)

	testVideoId := int64(10)
	testCommentCaseA := &model.Comment{
		VideoId: testVideoId,
		UserId:  5,
		Content: "test comment",
	}
	testCommentCaseB := &model.Comment{
		VideoId: testVideoId,
		UserId:  12,
		Content: "good",
	}
	testCommentCaseC := &model.Comment{
		VideoId: testVideoId,
		UserId:  19,
		Content: "hello",
	}
	db.Create(testCommentCaseA)
	db.Create(testCommentCaseB)
	db.Create(testCommentCaseC)

	expectedComments := []*model.Comment{
		testCommentCaseC,
		testCommentCaseB,
		testCommentCaseA,
	}

	c := NewCommentService(func() *gorm.DB {
		return db
	})
	comments, err := c.ListCommentOrderByCreatedAtDesc(testVideoId)
	assert.NoError(t, err)
	for i, comment := range comments {
		assert.Equal(t, expectedComments[i].VideoId, comment.VideoId)
		assert.Equal(t, expectedComments[i].UserId, comment.UserId)
		assert.Equal(t, expectedComments[i].Content, comment.Content)
	}
}

func TestCommentOperator_CountComment(t *testing.T) {
	db := newMockDB(t)

	testVideoId := int64(10)
	testCommentCaseA := &model.Comment{
		VideoId: testVideoId,
		UserId:  5,
		Content: "test comment",
	}
	testCommentCaseB := &model.Comment{
		VideoId: testVideoId,
		UserId:  12,
		Content: "good",
	}
	testCommentCaseC := &model.Comment{
		VideoId: testVideoId,
		UserId:  19,
		Content: "hello",
	}
	db.Create(testCommentCaseA)
	db.Create(testCommentCaseB)
	db.Create(testCommentCaseC)

	c := NewCommentService(func() *gorm.DB {
		return db
	})
	count, err := c.CountComment(testVideoId)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	t.Run("no comment", func(t *testing.T) {
		count, err := c.CountComment(100)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
