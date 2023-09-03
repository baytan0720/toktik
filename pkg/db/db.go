package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"toktik/pkg/config"
)

var db *gorm.DB

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
		cfg.Get(config.KEY_MYSQL_HOST).(string),
		cfg.Get(config.KEY_MYSQL_PORT).(int),
		cfg.Get(config.KEY_MYSQL_USER).(string),
		cfg.Get(config.KEY_MYSQL_PASSWORD).(string),
		cfg.Get(config.KEY_MYSQL_DATABASE).(string),
	)
	newDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	close()
	db = newDB
	log.Printf("db changed: %v:%v\n", cfg.Get(config.KEY_MYSQL_HOST), cfg.Get(config.KEY_MYSQL_PORT))
	return nil
}

func Instance() *gorm.DB {
	if db == nil {
		log.Panic("db is nil")
	}
	return db
}

func close() {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func generateDsn(host string, port int, user string, password string, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}
