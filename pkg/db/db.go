package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

var db *gorm.DB

func Init(dsn string) {
	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	err = db.AutoMigrate(&model.Comment{})
	if err != nil {
		log.Panic(err)
	}
}

func Instance() *gorm.DB {
	if db == nil {
		log.Panic("db is nil")
	}
	return db
}
