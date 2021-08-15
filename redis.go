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

func initRedisClient(config RedisConfig) *RedisClient {
  if config.Host == "" || config.Disabled {
    return nil
  }

  client := new(RedisClient)

  if c := client.GetDefaultClient(); c != nil {
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

  resclients[DefaultRedisClientKey] = client

  return client
}

func (r *RedisClient) GetClient(key string) *RedisClient {
  if resclients == nil {
    return nil
  }

  return resclients[key]
}

func (r *RedisClient) GetDefaultClient() *RedisClient {
  return r.GetClient(DefaultRedisClientKey)
}
