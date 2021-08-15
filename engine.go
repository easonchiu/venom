package venom

import (
  "github.com/gin-gonic/gin"
)

type Engine struct {
  Config *Config
  Redis  *RedisClient
  Mongo  *MongoClient
  Gin    *gin.Engine
}

// 初始化server
func Init(config *Config) *Engine {
  ginMode := gin.DebugMode
  if config.Mode == ModeProduction {
    ginMode = gin.ReleaseMode
  }

  gin.SetMode(ginMode)

  g := gin.Default()

  rds := initRedisClient(config.Redis)
  mgo := initMongoClient(config.Mongo)

  if !config.Logger.Disabled {
    g.Use(LoggerMiddleware(config.Logger))
  }

  engine := &Engine{
    Config: config,
    Redis:  rds,
    Mongo:  mgo,
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
  return e.Gin.Run(e.Config.Address + ":" + e.Config.Port)
}
