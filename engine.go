package venom

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

type Engine struct {
  Config *Config
  Redis  *RedisClient
  Mongo  *MongoClient
  Qmgo   *QmgoClient
  Gin    *gin.Engine
}

// 初始化server
func Init(config *Config) *Engine {
  ginMode := gin.DebugMode
  if config.Mode == ProductionMode {
    ginMode = gin.ReleaseMode
  }

  gin.SetMode(ginMode)

  g := gin.Default()

  rds := initRedisClient(DefaultRedisClientKey, config.Redis)
  mgo := initMongoClient(DefaultMongoClientKey, config.Mongo)
  qmgo := initQmgoClient(DefaultQmgoClientKey, config.Qmgo)

  if config.RedisMap != nil {
    for k, c := range config.RedisMap {
      initRedisClient(k, c)
    }
  }

  if config.MongoMap != nil {
    for k, c := range config.MongoMap {
      initMongoClient(k, c)
    }
  }

  if config.QmgoMap != nil {
    for k, c := range config.QmgoMap {
      initQmgoClient(k, c)
    }
  }

  if !config.Logger.Disabled {
    g.Use(LoggerMiddleware(config.Logger))
  }

  engine := &Engine{
    Config: config,
    Redis:  rds,
    Mongo:  mgo,
    Qmgo:   qmgo,
    Gin:    g,
  }

  return engine
}

// 使用中间件
func (e *Engine) Use(middleware ...Handle) {
  if middleware == nil || len(middleware) == 0 {
    return
  }

  funcs := make([]gin.HandlerFunc, 0, len(middleware))
  for _, m := range middleware {
    funcs = append(funcs, func(gctx *gin.Context) {
      m(&Context{Config: e.Config, Redis: e.Redis, Mongo: e.Mongo, GinContext: gctx})
    })
  }

  e.Gin.Use(funcs...)
}

// 启动server
func (e *Engine) Start() error {
  defer func() {
    e.Mongo.CloseAll()
    e.Qmgo.CloseAll()
    e.Redis.CloseAll()
  }()

  fmt.Println("Ready start venom ...")
  return e.Gin.Run(e.Config.Address + ":" + e.Config.Port)
}
