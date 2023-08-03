package operator

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

func newMockDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:toktik.db?&mode=memory"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&model.Comment{})
	require.NoError(t, err)
	return db
}

func TestCommentOperator_CreateComment(t *testing.T) {
	db := newMockDB(t)

	o := NewCommentOperator(db)
	comment, err := o.CreateComment(1, 1, "test")
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

	o := NewCommentOperator(db)
	comment, err := o.GetComment(testCommentCase.Id)
	assert.NoError(t, err)
	assert.Equal(t, testCommentCase.VideoId, comment.VideoId)
	assert.Equal(t, testCommentCase.UserId, comment.UserId)
	assert.Equal(t, testCommentCase.Content, comment.Content)
}

func TestCommentOperator_DeleteComment(t *testing.T) {
	db := newMockDB(t)

	testCommentCase := &model.Comment{
		VideoId: 10,
		UserId:  10,
		Content: "test comment",
	}
	db.Create(testCommentCase)

	o := NewCommentOperator(db)
	err := o.DeleteComment(testCommentCase.Id)
	assert.NoError(t, err)
	err = db.Where("id = ?", testCommentCase.Id).First(&model.Comment{}).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
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

	o := NewCommentOperator(db)
	comments, err := o.ListCommentOrderByCreatedAtDesc(testVideoId)
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

	o := NewCommentOperator(db)
	count, err := o.CountComment(testVideoId)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}
