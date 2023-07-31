package venom

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

type RedisPlugin struct {
	client *redis.Client
	config *RedisConfig
}

type RedisConfig struct {
	Name     string
	Host     string
	Port     string
	Password string
	DB       int
}

const DefaultRedisClientName = "default"

// redis 的 client list
var redisClients = make(map[string]*redis.Client)

// init plugin
func InitRedisPlugin(config *RedisConfig) *RedisPlugin {
	return &RedisPlugin{config: config}
}

// 启动
func (plugin *RedisPlugin) OnStart(config *Config) {
	if plugin.config.Name == "" {
		plugin.config.Name = DefaultRedisClientName
	}

	// client重名
	if _, exist := redisClients[plugin.config.Name]; exist {
		panic(fmt.Errorf("redis client named %v already exists", plugin.config.Name))
	}

	port := "6379"

	if plugin.config.Port != "" {
		port = plugin.config.Port
	}

	if strings.Contains(plugin.config.Host, ":") {
		sp := strings.Split(plugin.config.Host, ":")
		if len(sp) != 2 {
			panic(fmt.Errorf("redis host value is unavailable: %v", plugin.config.Host))
		}
		plugin.config.Host = sp[0]
		port = sp[1]
	}

	client := redis.NewClient(&redis.Options{
		Addr:     plugin.config.Host + ":" + port,
		Password: plugin.config.Password,
		DB:       plugin.config.DB,
	})

	plugin.client = client
	redisClients[plugin.config.Name] = client

	fmt.Printf("[PLUGIN] Redis <%v> start ok...\n", plugin.config.Name)
}

// 卸载
func (plugin *RedisPlugin) OnDestroy(config *Config) {
	if plugin.config.Name != "" && plugin != nil {
		if _, exist := redisClients[plugin.config.Name]; exist {
			_ = plugin.client.Close()
			delete(redisClients, plugin.config.Name)
		}
	}
}

// 获取 redis 的 client
func (plugin *RedisPlugin) GetClient(name ...string) *redis.Client {
	k := DefaultRedisClientName

	if name != nil {
		k = name[0]
	}

	if plugin.config != nil && plugin.config.Name == k && plugin.client != nil {
		return plugin.client
	}

	if redisClients == nil {
		return nil
	}

	return redisClients[k]
}

func GetRedisClient(client ...string) *redis.Client {
	plugin := new(RedisPlugin)
	return plugin.GetClient(client...)
}

// 检验是否实现了plugin interface
var _ IPlugin = (*RedisPlugin)(nil)
