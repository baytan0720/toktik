package video

import (
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"toktik/pkg/db/model"
	"toktik/pkg/test/testutil"
)

func newdb(t *testing.T) *gorm.DB {
	db := testutil.Newdb()
	require.NoError(t, db.AutoMigrate(&model.Relation{}))
	return db
}

func TestVideoService_GetVideo(t *testing.T) {
	db := newdb(t)
	// Create a mock video
	TestVideo := &model.Video{Id: 123, UserId: 456, Title: "Test Video"}
	require.NoError(t, db.Create(TestVideo).Error)
	var video model.Video
	require.NoError(t, db.First(&video, TestVideo.Id).Error)

	// Compare the actual video with the mock video
	require.Equal(t, video.Id, TestVideo.Id)
	require.Equal(t, video.UserId, TestVideo.UserId)
	require.Equal(t, video.Title, TestVideo.Title)
}

func TestVideoService_ListVideo(t *testing.T) {
	db := newdb(t)
	// 创建一个视频服务实例
	videoService := &VideoService{dbInstance: db}
	// 准备测试数据
	testUserID := int64(123)
	TestVideo1 := &Video{Id: 456, UserId: testUserID, Title: "Test Video 1"}
	TestVideo2 := &Video{Id: 789, UserId: testUserID, Title: "Test Video 2"}
	require.NoError(t, db.Create(TestVideo1).Error)
	require.NoError(t, db.Create(TestVideo2).Error)
	videos, err := videoService.ListVideo(testUserID)
	require.NoError(t, err)
	// 检查获取的视频列表是否正确
	expectedVideos := []Video{*TestVideo1, *TestVideo2}
	require.Equal(t, expectedVideos, videos)
}
func TestVideoService_GetVideoCount(t *testing.T) {
	db := newdb(t)
	userId := int64(123)
	// 假设预期的视频数量为10
	expectedCount := int64(1)
	videoService := &VideoService{}
	// 调用GetVideoCount函数获取视频数量
	count, err := videoService.GetVideoCount(userId)
	if err != nil {
		t.Errorf("GetVideoCount returned error: %v", err)
	}
	// 检查返回的视频数量是否与预期一致
	if count != expectedCount {
		t.Errorf("Expected count to be %d, but got %d", expectedCount, count)
	}
}

func TestVideoService_PublishVideo(t *testing.T) {
	db := newdb(t)
	videoService := VideoService{}
	err := videoService.PublishVideo("test_videoName", "test_imageName", 123, "test_title")
	if err != nil {
		t.Errorf("PublishVideo returned error : %v", err)
	}
}
