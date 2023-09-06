package config

var Conf Config

type Config interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Watch(key string, handler func(cfg Config))
	watch()
}

// config key
const (
	KEY_RELEASE        = "release"
	KEY_MYSQL          = "mysql"
	KEY_MYSQL_HOST     = "mysql.host"
	KEY_MYSQL_PORT     = "mysql.port"
	KEY_MYSQL_USER     = "mysql.user"
	KEY_MYSQL_PASSWORD = "mysql.password"
	KEY_MYSQL_DATABASE = "mysql.database"

	KEY_MINIO_ENDPOINT     = "minio.endpoint"
	KEY_MINIO_EXPOSE       = "minio.expose"
	KEY_MINIO_ACCESS_KEY   = "minio.access_key"
	KEY_MINIO_SECRET_KEY   = "minio.secret_key"
	KEY_MINIO_VIDEO_BUCKET = "minio.video_bucket"
	KEY_MINIO_COVER_BUCKET = "minio.cover_bucket"
)
