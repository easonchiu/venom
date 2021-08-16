package venom

import (
  "github.com/sirupsen/logrus"
)

type Mode int

const (
  ModeProduction Mode = iota
  ModeDevelopment
)

type Config struct {
  Address       string
  Port          string
  Mode          Mode
  SuccessFormat func(code int, obj interface{}) interface{}
  ErrorFormat   func(code int, errCode string, errMessage string, obj ...interface{}) interface{}
  ErrorCodes    map[string]string
  Apollo        ApolloConfig
  Redis         RedisConfig
  Mongo         MongoConfig
  RedisMap      map[string]RedisConfig
  MongoMap      map[string]MongoConfig
  Logger        LoggerConfig
}

type ApolloConfig struct {
  IP        string
  ID        string
  Cluster   string
  Namespace string
  Disabled  bool
}

type RedisConfig struct {
  Host     string
  Port     string
  Password string
  DB       int
  Disabled bool
}

type MongoConfig struct {
  URI      string
  Disabled bool
}

type LoggerConfig struct {
  Filename   string
  MaxSize    int
  MaxBackups int
  Level      logrus.Level
  MaxAge     int
  Disabled   bool
}
