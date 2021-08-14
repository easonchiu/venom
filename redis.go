package venom

import (
  goredis "github.com/go-redis/redis/v8"
)

type redisClients map[string]*goredis.Client

var resclients redisClients = make(map[string]*goredis.Client)

func initRedisClient(config RedisConfig) *goredis.Client {
  if config.Host == "" || config.Disabled {
    return nil
  }

  if client := GetRedisClient(); client != nil {
    return client
  }

  port := "6379"

  if config.Port != "" {
    port = config.Port
  }

  rds := goredis.NewClient(&goredis.Options{
    Addr:     config.Host + ":" + port,
    Password: config.Password,
    DB:       config.DB,
  })

  resclients["default"] = rds

  return rds
}

func GetRedisClient() *goredis.Client {
  if resclients != nil && resclients["default"] != nil {
    return resclients["default"]
  }

  return nil
}