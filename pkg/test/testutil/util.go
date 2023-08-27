package testutil

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewMockDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:toktik.db?&mode=memory"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
