package config

var Conf Config

type Config interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	Watch(key string, handler func(cfg Config))
	watch()
}

func Set(key string, value interface{}) {
	Conf.Set(key, value)
}

func Get(key string) interface{} {
	return Conf.Get(key)
}

func GetString(key string) string {
	return Conf.GetString(key)
}

func GetInt(key string) int {
	return Conf.GetInt(key)
}

func Watch(key string, handler func(cfg Config)) {
	Conf.Watch(key, handler)
}

// config key
const (
	KEY_RELEASE      = "release"
	KEY_SERVICE_NAME = "name"
	KEY_LISTEN_ON    = "listen_on"

	KEY_LOGGER_LEVEL  = "logger.level"
	KEY_LOGGER_OUTPUT = "logger.output"

	KEY_CONSUL = "consul"

	KEY_MYSQL          = "mysql"
	KEY_MYSQL_HOST     = "mysql.host"
	KEY_MYSQL_PORT     = "mysql.port"
	KEY_MYSQL_USER     = "mysql.user"
	KEY_MYSQL_PASSWORD = "mysql.password"
	KEY_MYSQL_DATABASE = "mysql.database"

	KEY_MINIO_ENDPOINT          = "minio.endpoint"
	KEY_MINIO_EXPOSE            = "minio.expose"
	KEY_MINIO_ACCESS_KEY        = "minio.access_key"
	KEY_MINIO_SECRET_KEY        = "minio.secret_key"
	KEY_MINIO_VIDEO_BUCKET      = "minio.video_bucket"
	KEY_MINIO_COVER_BUCKET      = "minio.cover_bucket"
	KEY_MINIO_AVATAR_BUCKET     = "minio.avatar_bucket"
	KEY_MINIO_BACKGRAUND_BUCKET = "minio.background_bucket"
	KEY_MINIO_LOCATION          = "minio.location"

	KEY_RABBITMQ             = "rabbitmq"
	KEY_RABBITMQ_HOST        = "rabbitmq.host"
	KEY_RABBITMQ_PORT        = "rabbitmq.port"
	KEY_RABBITMQ_USER        = "rabbitmq.user"
	KEY_RABBITMQ_PASSWORD    = "rabbitmq.password"
	KEY_RABBITMQ_QUEUE       = "rabbitmq.queue"
	KEY_RABBITMQ_EXCHANGE    = "rabbitmq.exchange"
	KEY_RABBITMQ_ROUTING_KEY = "rabbitmq.routing_key"
)
