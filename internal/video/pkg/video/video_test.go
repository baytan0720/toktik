package video

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
	require.NoError(t, db.AutoMigrate(&model.Video{}))
	return db
}

func TestVideoOperator_GetFeed(t *testing.T) {
	db := newMockDB(t)

	r := NewVideoService(func() *gorm.DB {
		return db
	})

	testVideoCaseA := &model.Video{
		UserId:    10,
		CreatedAt: time.Now().Add(-time.Hour * 24 * 5),
	}

	testVideoCaseB := &model.Video{
		UserId:    10,
		CreatedAt: time.Now().Add(-time.Hour * 24 * 3),
	}

	testVideoCaseC := &model.Video{
		UserId:    11,
		CreatedAt: time.Now().Add(-time.Hour * 24 * 1),
	}
	db.Create(testVideoCaseA)
	db.Create(testVideoCaseB)
	db.Create(testVideoCaseC)

	videosBeforeTwoHoursAgo, err := r.GetFeed(time.Now().Add(-time.Hour * 24 * 2).Unix() * time.Microsecond.Nanoseconds())
	require.NoError(t, err)
	assert.Equal(t, 2, len(videosBeforeTwoHoursAgo))

	videosBeforeFourHoursAgo, err := r.GetFeed(time.Now().Add(-time.Hour * 24 * 4).Unix() * time.Microsecond.Nanoseconds())
	require.NoError(t, err)
	assert.Equal(t, 1, len(videosBeforeFourHoursAgo))

	videosAllTime, err := r.GetFeed(int64(0))
	require.NoError(t, err)
	assert.Equal(t, 3, len(videosAllTime))

}
