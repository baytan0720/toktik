package relation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/pkg/db/model"
	"toktik/pkg/test/testutil"
)

func newMockDB(t *testing.T) *gorm.DB {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&model.Relation{}))
	return db
}

func TestRelationOperator_GetFollowInfo(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	testRelationCaseA := &model.Relation{
		UserId:    10,
		ToUserId:  11,
		IsFollow:  true,
	}

	testRelationCaseB := &model.Relation{
		UserId:    10,
		ToUserId:  12,
		IsFollow:  false,
	}

	testRelationCaseC := &model.Relation{
		UserId:    11,
		ToUserId:  12,
		IsFollow:  true,
	}

	db.Create(testRelationCaseA)
	db.Create(testRelationCaseB)
	db.Create(testRelationCaseC)

	// Test GetFollowRelations
	relations, err := r.GetFollowRelations(10, []int64{11, 12, 13})
	require.NoError(t, err)
	assert.Equal(t, 2, len(relations))

	// Test GetFollowCount
	count, err := r.GetFollowCount(10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Test GetFollowerCount
	count, err = r.GetFollowerCount(12)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestRelationOperator_Follow(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	relation := Relation{
		UserId: 0, ToUserId: 10,
	}
	relation2 := Relation{
		UserId: 10, ToUserId: 20,
	}

	testRelationCase := &model.Relation{
		UserId:    0,
		ToUserId:  10,
		IsFollow:  false,
	}

	db.Create(testRelationCase)

	// Test upsert
	err := r.Follow(0, 10)
	require.NoError(t, err)
	db.First(&relation)
	assert.True(t, relation.IsFollow)

	// Test create
	err = r.Follow(10, 20)
	require.NoError(t, err)
	db.First(&relation2)
	assert.Equal(t, true, relation.IsFollow)
}

func TestRelationOperator_UnFollow(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	relation := Relation{
		UserId: 1, ToUserId: 10,
	}

	testRelationCase := &model.Relation{
		UserId:    1,
		ToUserId:  10,
		IsFollow:  true,
	}

	db.Create(testRelationCase)

	// Test Unfollow
	require.NoError(t, r.Unfollow(1, 10))
	db.First(&relation)
	assert.False(t, relation.IsFollow)

}

func TestRelationOperator_IsFollow(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	testRelationCase := &model.Relation{
		UserId:    1,
		ToUserId:  10,
		IsFollow:  true,
	}

	db.Create(testRelationCase)

	// 存在记录
	isFollow, err := r.IsFollow(1, 10)
	require.NoError(t, err)
	assert.True(t, isFollow)

	// 不存在记录
	isFollow, err = r.IsFollow(2, 10)
	require.NoError(t, err)
	assert.False(t, isFollow)
}

func TestRelationOperator_ListFollow(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	testRelationCaseA := &model.Relation{
		UserId:    10,
		ToUserId:  11,
		IsFollow:  true,
		UpdatedAt: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	testRelationCaseB := &model.Relation{
		UserId:    10,
		ToUserId:  12,
		IsFollow:  true,
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	testRelationCaseC := &model.Relation{
		UserId:    10,
		ToUserId:  13,
		IsFollow:  true,
		UpdatedAt: time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	db.Create(testRelationCaseA)
	db.Create(testRelationCaseB)
	db.Create(testRelationCaseC)

	// Test ListFollow
	followList, err := r.ListFollow(10)
	require.NoError(t, err)
	assert.Equal(t, 3, len(followList))
	assert.Equal(t, int64(13), followList[0])
	assert.Equal(t, int64(11), followList[1])
	assert.Equal(t, int64(12), followList[2])
}

func TestRelationOperator_ListFollower(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	testRelationCaseA := &model.Relation{
		UserId:    11,
		ToUserId:  10,
		IsFollow:  true,
		UpdatedAt: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	testRelationCaseB := &model.Relation{
		UserId:    12,
		ToUserId:  10,
		IsFollow:  true,
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	testRelationCaseC := &model.Relation{
		UserId:    13,
		ToUserId:  10,
		IsFollow:  true,
		UpdatedAt: time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	db.Create(testRelationCaseA)
	db.Create(testRelationCaseB)
	db.Create(testRelationCaseC)

	// Test ListFollower
	followList, err := r.ListFollower(10)
	require.NoError(t, err)
	assert.Equal(t, 3, len(followList))
	assert.Equal(t, int64(13), followList[0])
	assert.Equal(t, int64(11), followList[1])
	assert.Equal(t, int64(12), followList[2])
}