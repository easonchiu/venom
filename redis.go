package venom

import (
  "github.com/go-redis/redis/v8"
)

type RedisClient struct {
  *redis.Client
}

const DefaultRedisClientKey = "default"

var rdsclients = make(map[string]*RedisClient)

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

  rdsclients[key] = client

  return client
}

func (r *RedisClient) GetClient(key ...string) *RedisClient {
  k := DefaultRedisClientKey

  if key != nil && len(key) > 0 {
    k = key[0]
  }

  if rdsclients == nil {
    return nil
  }

  return rdsclients[k]
}

func (r *RedisClient) GetDefaultClient() *RedisClient {
  return r.GetClient(DefaultRedisClientKey)
}

func (r *RedisClient) CloseAll() {
  for k, c := range rdsclients {
    if c != nil {
      _ = c.Close()
      rdsclients[k] = nil
    }
  }
}