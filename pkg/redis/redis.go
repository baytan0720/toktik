package redis

import (
	"github.com/go-redis/redis/v7"

	"toktik/pkg/config"
)

type RedisBucket int

var redisClient *redis.Client

const (
	UserBucket RedisBucket = iota
	VideoBucket
	CommentBucket
	FavoriteBucket
	RelationBucket
)

func InitRedisClient() {
	var bucket RedisBucket
	switch config.GetString(config.KEY_SERVICE_NAME) {
	case "user":
		bucket = UserBucket
	case "video":
		bucket = VideoBucket
	case "comment":
		bucket = CommentBucket
	case "favorite":
		bucket = FavoriteBucket
	case "relation":
		bucket = RelationBucket
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.addr"),
		Password: config.GetString("redis.password"),
		DB:       int(bucket),
	})
}

func Instance() *redis.Client {
	if redisClient == nil {
		InitRedisClient()
	}
	return redisClient
}
