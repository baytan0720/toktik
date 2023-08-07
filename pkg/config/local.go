package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type LocalConf struct {
	conf     *viper.Viper
	handlers map[string]func(cfg Config)
}

var _ Config = &LocalConf{}

func ReadConfigFromLocal(configPath string) Config {
	conf := viper.New()
	conf.SetConfigFile(configPath)
	conf.SetConfigType("yaml")
	err := conf.ReadInConfig()
	if err != nil {
		log.Fatalln("Read config file failed: ", err)
	}

	c := &LocalConf{
		conf: conf,
	}
	Conf = c
	go c.watch()
	return Conf
}

func (c *LocalConf) Set(key string, value interface{}) {
	c.conf.Set(key, value)
}

func (c *LocalConf) Get(key string) interface{} {
	return c.conf.Get(key)
}

func (c *LocalConf) Watch(key string, handler func(cfg Config)) {
	c.conf.WatchConfig()
}

func (c *LocalConf) watch() {
	c.conf.OnConfigChange(func(in fsnotify.Event) {
		for _, handler := range c.handlers {
			handler(c)
		}
	})
}
