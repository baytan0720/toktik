package db

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"toktik/pkg/config"
	"toktik/pkg/db/model"
)

var db *gorm.DB
var mu sync.Mutex

func Init() {
	cfg := config.Conf
	err := Use(cfg)
	if err != nil {
		log.Fatalln("failed to connect database: ", err)
	}
	cfg.Watch(config.KEY_MYSQL, func(cfg config.Config) {
		err = Use(cfg)
		if err != nil {
			log.Println("failed to connect database: ", err)
		}
	})
}

func Use(cfg config.Config) error {
	dsn := generateDsn(
		cfg.GetString(config.KEY_MYSQL_HOST),
		cfg.GetInt(config.KEY_MYSQL_PORT),
		cfg.GetString(config.KEY_MYSQL_USER),
		cfg.GetString(config.KEY_MYSQL_PASSWORD),
		cfg.GetString(config.KEY_MYSQL_DATABASE),
	)
	newDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	mu.Lock()
	db = newDB
	mu.Unlock()
	err = db.AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{}, &model.Relation{})
	if err != nil {
		return err
	}
	log.Printf("db changed: %v:%v\n", cfg.Get(config.KEY_MYSQL_HOST), cfg.Get(config.KEY_MYSQL_PORT))
	return nil
}

func Instance() *gorm.DB {
	mu.Lock()
	defer mu.Unlock()
	if db == nil {
		log.Panic("db is nil")
	}
	return db
}

func generateDsn(host string, port int, user string, password string, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}
