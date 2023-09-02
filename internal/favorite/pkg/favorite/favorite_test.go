package favorite

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"time"
	"toktik/pkg/db/model"
	"toktik/pkg/test/testutil"
)

func newMockDB(t *testing.T) *gorm.DB {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&model.Favorite{}))
	return db
}

func TestFavoriteService_Favorite(t *testing.T) {
	db := newMockDB(t)

	f := NewFavoriteService(func() *gorm.DB {
		return db
	})

	favorite1 := Favorite{
		UserId: 0, VideoId: 10,
	}
	favorite2 := Favorite{
		UserId: 10, VideoId: 20,
	}

	testFavoriteCase := &model.Favorite{
		UserId:     0,
		VideoId:    10,
		IsFavorite: false,
	}

	db.Create(testFavoriteCase)

	// Test upsert
	err := f.Favorite(10, 0)
	require.NoError(t, err)
	db.First(&favorite1)
	assert.True(t, favorite1.IsFavorite)

	// Test create
	err = f.Favorite(20, 10)
	require.NoError(t, err)
	db.First(&favorite2)
	assert.True(t, favorite2.IsFavorite)
}

func TestFavoriteService_UnFavorite(t *testing.T) {
	db := newMockDB(t)

	testFavoriteCase := &model.Favorite{
		VideoId:    10,
		UserId:     10,
		IsFavorite: true,
	}
	db.Create(testFavoriteCase)

	f := NewFavoriteService(func() *gorm.DB {
		return db
	})
	err := f.UnFavorite(testFavoriteCase.VideoId, testFavoriteCase.UserId)
	assert.NoError(t, err)
	err = db.Where("video_id=? AND user_id=? AND is_favorite", testFavoriteCase.VideoId, testFavoriteCase.UserId, true).First(&Favorite{}).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestFavoriteService_ListFavorite(t *testing.T) {
	db := newMockDB(t)
	testUserId := int64(10)
	testFavoriteCaseA := &model.Favorite{
		VideoId:    3,
		UserId:     testUserId,
		IsFavorite: true,
		UpdatedAt:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	testFavoriteCaseB := &model.Favorite{
		VideoId:    12,
		UserId:     testUserId,
		IsFavorite: true,
		UpdatedAt:  time.Date(2001, 1, 1, 0, 1, 0, 0, time.UTC),
	}
	testFavoriteCaseC := &model.Favorite{
		VideoId:    8,
		UserId:     testUserId,
		IsFavorite: true,
		UpdatedAt:  time.Date(2001, 1, 1, 1, 0, 0, 0, time.UTC),
	}
	db.Create(testFavoriteCaseA)
	db.Create(testFavoriteCaseB)
	db.Create(testFavoriteCaseC)

	expectedFavorite := []*model.Favorite{
		testFavoriteCaseC,
		testFavoriteCaseB,
		testFavoriteCaseA,
	}

	f := NewFavoriteService(func() *gorm.DB {
		return db
	})
	videos, err := f.ListFavorite(testUserId)
	assert.NoError(t, err)
	for i, videoId := range videos {
		assert.Equal(t, expectedFavorite[i].VideoId, videoId)
	}
}

func TestFavoriteService_CountVideoFavorite(t *testing.T) {
	db := newMockDB(t)

	testVideoId := int64(10)
	testFavoriteCaseA := &model.Favorite{
		VideoId:    testVideoId,
		UserId:     13,
		IsFavorite: true,
	}
	testFavoriteCaseB := &model.Favorite{
		VideoId:    testVideoId,
		UserId:     7,
		IsFavorite: true,
	}
	testFavoriteCaseC := &model.Favorite{
		VideoId:    testVideoId,
		UserId:     5,
		IsFavorite: true,
	}
	db.Create(testFavoriteCaseA)
	db.Create(testFavoriteCaseB)
	db.Create(testFavoriteCaseC)

	f := NewFavoriteService(func() *gorm.DB {
		return db
	})
	count, err := f.CountVideoFavorite(testVideoId)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestFavoriteService_CountUserFavorite(t *testing.T) {
	db := newMockDB(t)

	testUserId := int64(10)
	testFavoriteCaseA := &model.Favorite{
		UserId:     testUserId,
		VideoId:    13,
		IsFavorite: true,
	}
	testFavoriteCaseB := &model.Favorite{
		UserId:     testUserId,
		VideoId:    7,
		IsFavorite: true,
	}
	testFavoriteCaseC := &model.Favorite{
		UserId:     testUserId,
		VideoId:    5,
		IsFavorite: false,
	}
	db.Create(testFavoriteCaseA)
	db.Create(testFavoriteCaseB)
	db.Create(testFavoriteCaseC)

	f := NewFavoriteService(func() *gorm.DB {
		return db
	})
	count, err := f.CountUserFavorite(testUserId)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestFavoriteService_IsFavorite(t *testing.T) {
	db := newMockDB(t)

	testFavoriteCase := &model.Favorite{
		VideoId:    5,
		UserId:     25,
		IsFavorite: true,
	}
	db.Create(testFavoriteCase)

	f := NewFavoriteService(func() *gorm.DB {
		return db
	})
	isFavorite, err := f.IsFavorite(25, 5)
	assert.NoError(t, err)
	assert.True(t, isFavorite)

	isFavorite, err = f.IsFavorite(25, 4)
	assert.NoError(t, err)
	assert.False(t, isFavorite)
}

