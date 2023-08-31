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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testRelationCaseB := &model.Relation{
		UserId:    10,
		ToUserId:  12,
		IsFollow:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testRelationCaseC := &model.Relation{
		UserId:    11,
		ToUserId:  12,
		IsFollow:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
		UserId: 0, ToUserId: 10,
	}

	testRelationCase := &model.Relation{
		UserId:    0,
		ToUserId:  10,
		IsFollow:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.Create(testRelationCase)

	err := r.Unfollow(0, 10)
	db.First(&relation)
	require.NoError(t, err)
	assert.False(t, relation.IsFollow)
}

func TestRelationOperator_GetFollow(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	testRelationCaseA := &model.Relation{
		UserId:    10,
		ToUserId:  11,
		IsFollow:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testRelationCaseB := &model.Relation{
		UserId:    10,
		ToUserId:  12,
		IsFollow:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testRelationCaseC := &model.Relation{
		UserId:    10,
		ToUserId:  13,
		IsFollow:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.Create(testRelationCaseA)
	db.Create(testRelationCaseB)
	db.Create(testRelationCaseC)

	// Test GetFollow
	followList, err := r.GetFollow(10)
	require.NoError(t, err)
	assert.Equal(t, 2, len(followList))
	assert.Equal(t, int64(13), followList[1])
}

func TestRelationOperator_GetFollower(t *testing.T) {
	db := newMockDB(t)

	r := NewRelationService(func() *gorm.DB {
		return db
	})

	testRelationCaseA := &model.Relation{
		UserId:    10,
		ToUserId:  12,
		IsFollow:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testRelationCaseB := &model.Relation{
		UserId:    11,
		ToUserId:  12,
		IsFollow:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testRelationCaseC := &model.Relation{
		UserId:    13,
		ToUserId:  12,
		IsFollow:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.Create(testRelationCaseA)
	db.Create(testRelationCaseB)
	db.Create(testRelationCaseC)

	// Test GetFollower
	followList, err := r.GetFollower(12)
	require.NoError(t, err)
	assert.Equal(t, 2, len(followList))
	assert.Equal(t, int64(11), followList[0])
}