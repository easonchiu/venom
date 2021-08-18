package venom

import (
  "github.com/sirupsen/logrus"
)

type Mode int

const (
  ProductionMode Mode = iota
  DevelopmentMode
)

type Config struct {
  Address       string
  Port          string
  Mode          Mode
  SuccessFormat func(code int, obj interface{}) interface{}
  FailFormat    func(code int, errCode string, errMessage string, obj ...interface{}) interface{}
  FailCodes     map[string]string
  Apollo        ApolloConfig
  Redis         RedisConfig
  Mongo         MongoConfig
  Qmgo          QmgoConfig
  RedisMap      map[string]RedisConfig
  MongoMap      map[string]MongoConfig
  QmgoMap       map[string]QmgoConfig
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
  URI         string
  Database    string
  MinPoolSize uint64
  MaxPoolSize uint64
  Disabled    bool
}

type QmgoConfig struct {
  URI         string
  Database    string
  MinPoolSize uint64
  MaxPoolSize uint64
  Disabled    bool
}

type LoggerConfig struct {
  Filename   string
  MaxSize    int
  MaxBackups int
  Level      logrus.Level
  MaxAge     int
  Disabled   bool
}
