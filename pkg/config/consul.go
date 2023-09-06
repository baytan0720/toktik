package config

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

// ConsulConf is a config struct for consul
type ConsulConf struct {
	conf         sync.Map
	consulClient *consulapi.Client
	decoder      map[string]func(cfg Config, pair *consulapi.KVPair)
	handlers     map[string]func(cfg Config)
}

// check interface
var _ Config = &ConsulConf{}

// ReadConfigFromConsul read config from consul
func ReadConfigFromConsul(consulAddr string) Config {
	// connect to consul
	consulClient, err := consulapi.NewClient(&consulapi.Config{Address: consulAddr})
	if err != nil {
		log.Fatalf("Connect to consul failed: %v", err)
	}

	// get all key-value pairs
	pairs, _, err := consulClient.KV().List("", nil)
	if err != nil {
		log.Fatalf("Get consul kv failed: %v", err)
	}

	c := &ConsulConf{
		conf:         sync.Map{},
		consulClient: consulClient,
		decoder:      newDecoder(),
		handlers:     make(map[string]func(cfg Config)),
	}

	// set key-value pairs to config
	for _, pair := range pairs {
		c.Set(pair.Key, string(pair.Value))
		if decoder := c.decoder[pair.Key]; decoder != nil {
			decoder(c, pair)
		}
	}

	Conf = c
	go c.watch()
	return Conf
}

// Set key-value to config
func (c *ConsulConf) Set(key string, value interface{}) {
	c.conf.Store(key, value)
}

// Get value from config
func (c *ConsulConf) Get(key string) interface{} {
	val, _ := c.conf.Load(key)
	return val
}

// GetString value from config
func (c *ConsulConf) GetString(key string) string {
	val, _ := c.conf.Load(key)
	return val.(string)
}

// GetInt value from config
func (c *ConsulConf) GetInt(key string) int {
	val, _ := c.conf.Load(key)
	return val.(int)
}

// Watch config will update key-value and call handler when config changed
func (c *ConsulConf) Watch(key string, handler func(cfg Config)) {
	if handler == nil {
		return
	}
	c.handlers[key] = handler
}

func (c *ConsulConf) watch() {
	interval := 10 * time.Second

	for {
		time.Sleep(interval)

		var wg sync.WaitGroup
		for key, handler := range c.handlers {
			wg.Add(1)

			go func(key string, handler func(cfg Config)) {
				defer wg.Done()
				pair, _, err := c.consulClient.KV().Get(key, nil)
				if err != nil {
					log.Printf("Get %s from consul failed: %v", key, err)
					return
				}
				oldValue, ok := c.Get(pair.Key).(string)
				if !ok {
					log.Printf("Unexpected type of %s: %v", pair.Key, c.Get(pair.Key))
				}
				newValue := pair.Value
				if string(newValue) != oldValue && len(newValue) > 0 {
					c.Set(pair.Key, string(newValue))
					if decoder := c.decoder[pair.Key]; decoder != nil {
						decoder(c, pair)
					}
					handler(c)
				}
			}(key, handler)
		}

		wg.Wait()
	}
}

func newDecoder() map[string]func(cfg Config, pair *consulapi.KVPair) {
	return map[string]func(cfg Config, pair *consulapi.KVPair){
		KEY_MYSQL: func(cfg Config, pair *consulapi.KVPair) {
			payload := make(map[string]interface{})
			err := json.Unmarshal(pair.Value, &payload)
			if err != nil {
				log.Printf("Unmarshal mysql config failed: %v", err)
				return
			}
			for key, value := range payload {
				if key == "port" {
					port, ok := value.(float64)
					if !ok {
						log.Printf("Unexpected type of mysql.port: %v", value)
						continue
					}
					value = int(port)
				}
				cfg.Set("mysql."+key, value)
			}
		},
	}
}
