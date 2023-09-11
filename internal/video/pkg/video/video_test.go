package video

import (
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/pkg/test/testutil"
)

func newMockDB(t *testing.T) *gorm.DB {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&Video{}))
	return db
}

func TestVideoOperator_Create(t *testing.T) {
	db := newMockDB(t)
	s := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	v := NewVideoService(func() *gorm.DB {
		return db
	}, func() *redis.Client {
		return rdb
	})

	err := v.CreateVideo(1, "title1", "", "")
	require.NoError(t, err)

	err = v.CreateVideo(2, "title2", "", "")
	require.NoError(t, err)

	var videos []Video
	db.Find(&videos)
	assert.Equal(t, 2, len(videos))

	keys := []string{strconv.FormatInt(videos[0].Id, 10), strconv.FormatInt(videos[1].Id, 10)}
	video1, err := rdb.HGetAll(keys[0]).Result()
	video2, err := rdb.HGetAll(keys[1]).Result()
	require.NoError(t, err)

	userId, _ := strconv.ParseInt(video1["UserId"], 10, 64)
	assert.Equal(t, int64(1), userId)
	assert.Equal(t, "title1", video1["Title"])

	userId, _ = strconv.ParseInt(video2["UserId"], 10, 64)
	assert.Equal(t, int64(2), userId)
	assert.Equal(t, "title2", video2["Title"])
}

func TestVideoService_ListVideoByUserId(t *testing.T) {
	db := newMockDB(t)
	s := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	v := NewVideoService(func() *gorm.DB {
		return db
	}, func() *redis.Client {
		return rdb
	})

	rdb.HSet("1", map[string]interface{}{
		"Id":       1,
		"UserId":   1,
		"Title":    "title1",
		"PlayUrl":  "",
		"CoverUrl": "",
		"Existed":  true,
	})
	rdb.HSet("2", map[string]interface{}{
		"Id":       2,
		"UserId":   1,
		"Title":    "title2",
		"PlayUrl":  "",
		"CoverUrl": "",
		"Existed":  true,
	})
	rdb.HSet("3", map[string]interface{}{
		"Id":       3,
		"UserId":   2,
		"Title":    "title3",
		"PlayUrl":  "",
		"CoverUrl": "",
		"Existed":  true,
	})
	rdb.SAdd("user_videos:1", "1", "2")
	rdb.SAdd("user_videos:2", "3")

	videos, err := v.ListVideoByUserId(1)
	require.NoError(t, err)
	assert.Equal(t, 2, len(videos))
	assert.Equal(t, int64(1), videos[0].Id)
	assert.Equal(t, int64(1), videos[0].UserId)
	assert.Equal(t, int64(2), videos[1].Id)
	assert.Equal(t, int64(1), videos[1].UserId)

	videos, err = v.ListVideoByUserId(2)
	require.NoError(t, err)
	assert.Equal(t, 1, len(videos))
	assert.Equal(t, int64(3), videos[0].Id)
	assert.Equal(t, int64(2), videos[0].UserId)

	videos, err = v.ListVideoByUserId(3)
	require.NoError(t, err)
	assert.Equal(t, 0, len(videos))

}

func TestVideoOperator_GetVideoByIds(t *testing.T) {
	db := newMockDB(t)
	s := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	v := NewVideoService(func() *gorm.DB {
		return db
	}, func() *redis.Client {
		return rdb
	})

	rdb.HSet("1", map[string]interface{}{
		"Id":       1,
		"UserId":   1,
		"Title":    "title1",
		"PlayUrl":  "",
		"CoverUrl": "",
		"Existed":  true,
	})
	rdb.HSet("2", map[string]interface{}{
		"Id":       2,
		"UserId":   1,
		"Title":    "title2",
		"PlayUrl":  "",
		"CoverUrl": "",
		"Existed":  true,
	})
	rdb.HSet("3", map[string]interface{}{
		"Id":       3,
		"UserId":   2,
		"Title":    "title3",
		"PlayUrl":  "",
		"CoverUrl": "",
		"Existed":  true,
	})
	db.Create(&Video{
		Id: 4,
		UserId: 3,
		Title: "title4",
		PlayUrl: "",
		CoverUrl: "",
	})

	videos, err := v.GetVideoByIds([]int64{1, 2, 3})
	require.NoError(t, err)
	assert.Equal(t, 3, len(videos))
	assert.Equal(t, int64(1), videos[0].Id)
	assert.Equal(t, int64(1), videos[0].UserId)
	assert.Equal(t, int64(2), videos[1].Id)
	assert.Equal(t, int64(1), videos[1].UserId)
	assert.Equal(t, int64(3), videos[2].Id)
	assert.Equal(t, int64(2), videos[2].UserId)

	videos, err = v.GetVideoByIds([]int64{1, 2})
	require.NoError(t, err)
	assert.Equal(t, 2, len(videos))
	assert.Equal(t, int64(1), videos[0].Id)
	assert.Equal(t, int64(1), videos[0].UserId)
	assert.Equal(t, int64(2), videos[1].Id)
	assert.Equal(t, int64(1), videos[1].UserId)

	videos, err = v.GetVideoByIds([]int64{4})
	require.NoError(t, err)
	assert.Equal(t, int64(4), videos[0].Id)
	assert.Equal(t, int64(3), videos[0].UserId)
	assert.Equal(t, "title4", videos[0].Title)

	videos, err = v.GetVideoByIds([]int64{5})
	assert.ErrorContains(t, err, "video not found")
}

func TestVideoOperator_CountWork(t *testing.T) {
	db := newMockDB(t)
	s := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	v := NewVideoService(func() *gorm.DB {
		return db
	}, func() *redis.Client {
		return rdb
	})

	testVideoCaseA := &Video{
		UserId: 10,
	}

	testVideoCaseB := &Video{
		UserId: 10,
	}

	testVideoCaseC := &Video{
		UserId: 11,
	}
	db.Create(testVideoCaseA)
	db.Create(testVideoCaseB)
	db.Create(testVideoCaseC)

	count, err := v.CountWork(10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	count, err = v.CountWork(11)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = v.CountWork(12)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestVideoOperator_GetFeed(t *testing.T) {
	db := newMockDB(t)
	s := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	v := NewVideoService(func() *gorm.DB {
		return db
	}, func() *redis.Client {
		return rdb
	})
	testVideoCaseA := &Video{
		UserId:    10,
		CreatedAt: time.Now().Add(-time.Hour * 24 * 5),
	}

	testVideoCaseB := &Video{
		UserId:    10,
		CreatedAt: time.Now().Add(-time.Hour * 24 * 3),
	}

	testVideoCaseC := &Video{
		UserId:    11,
		CreatedAt: time.Now().Add(-time.Hour * 24 * 1),
	}
	db.Create(testVideoCaseA)
	db.Create(testVideoCaseB)
	db.Create(testVideoCaseC)

	videosBeforeTwoHoursAgo, err := v.GetFeed(time.Now().Add(-time.Hour*24*2).Unix() * time.Microsecond.Nanoseconds())
	require.NoError(t, err)
	assert.Equal(t, 2, len(videosBeforeTwoHoursAgo))

	videosBeforeFourHoursAgo, err := v.GetFeed(time.Now().Add(-time.Hour*24*4).Unix() * time.Microsecond.Nanoseconds())
	require.NoError(t, err)
	assert.Equal(t, 1, len(videosBeforeFourHoursAgo))

	videosAllTime, err := v.GetFeed(int64(0))
	require.NoError(t, err)
	assert.Equal(t, 3, len(videosAllTime))

}
