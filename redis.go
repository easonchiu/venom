package venom

import (
  "github.com/go-redis/redis/v8"
)

type RedisClient struct {
  *redis.Client
}

type redisClients map[string]*RedisClient

const DefaultRedisClientKey = "default"

var resclients redisClients = make(map[string]*RedisClient)

func initRedisClient(key string, config RedisConfig) *RedisClient {
  client := new(RedisClient)

  if key == "" || config.Host == "" || config.Disabled {
    return client
  }

  if c := client.GetClient(key); c != nil {
    return c
  }

  port := "6379"

  if config.Port != "" {
    port = config.Port
  }

  client = &RedisClient{
    redis.NewClient(&redis.Options{
      Addr:     config.Host + ":" + port,
      Password: config.Password,
      DB:       config.DB,
    }),
  }

  resclients[key] = client

  return client
}

func (r *RedisClient) GetClient(key ...string) *RedisClient {
  k := DefaultRedisClientKey

  if key != nil && len(key) > 0 {
    k = key[0]
  }

  if resclients == nil {
    return nil
  }

  return resclients[k]
}

func (r *RedisClient) GetDefaultClient() *RedisClient {
  return r.GetClient(DefaultRedisClientKey)
}
