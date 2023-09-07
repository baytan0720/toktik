package user

import (
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"toktik/pkg/db/model"
	"toktik/pkg/test/testutil"
)

func newMockDB(t *testing.T) *gorm.DB {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&model.User{}))
	return db
}

func TestUserService_CreateUser(t *testing.T) {
	db := newMockDB(t)
	s := NewUserService(func() *gorm.DB {
		return db
	})

	t.Run("success", func(t *testing.T) {
		_, err := s.CreateUser("test", "123456")
		assert.NoError(t, err)
	})

	t.Run("duplicate", func(t *testing.T) {
		_, err := s.CreateUser("test", "123456")
		assert.Error(t, err)
	})
}

func TestUserService_Login(t *testing.T) {
	db := newMockDB(t)
	s := NewUserService(func() *gorm.DB {
		return db
	})

	t.Run("username not found", func(t *testing.T) {
		_, err := s.Login("test", "123456")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("success", func(t *testing.T) {
		err := db.Create(&model.User{
			Username: "test",
			Password: fmt.Sprintf("%x", md5.Sum([]byte("123456"))),
		}).Error
		require.NoError(t, err)
		_, err = s.Login("test", "123456")
		assert.NoError(t, err)
	})

	t.Run("password incorrect", func(t *testing.T) {
		_, err := s.Login("test", "1234567")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestUserService_GetUserById(t *testing.T) {
	db := newMockDB(t)
	s := NewUserService(func() *gorm.DB {
		return db
	})

	err := db.Create(&model.User{
		Username: "test",
		Password: fmt.Sprintf("%x", md5.Sum([]byte("123456"))),
	}).Error
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		user, err := s.GetUserById(1)
		assert.NoError(t, err)
		assert.Equal(t, "test", user.Username)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := s.GetUserById(2)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestUserService_GetUserByIds(t *testing.T) {
	db := newMockDB(t)
	s := NewUserService(func() *gorm.DB {
		return db
	})

	err := db.Create(&model.User{
		Username: "test",
		Password: fmt.Sprintf("%x", md5.Sum([]byte("123456"))),
	}).Error
	require.NoError(t, err)

	err = db.Create(&model.User{
		Username: "test2",
		Password: fmt.Sprintf("%x", md5.Sum([]byte("123456"))),
	}).Error
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		users, err := s.GetUserByIds([]int64{1, 2})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(users))
	})

	t.Run("1 not found", func(t *testing.T) {
		users, err := s.GetUserByIds([]int64{1, 3})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
}
