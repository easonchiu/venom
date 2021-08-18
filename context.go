package venom

import (
  "context"
  "github.com/gin-gonic/gin"
)

type Context struct {
  Config     *Config
  Redis      *RedisClient
  Mongo      *MongoClient
  Qmgo       *QmgoClient
  GinContext *gin.Context
}

func (ctx *Context) Bg() context.Context {
  return context.Background()
}

func (ctx *Context) Success(code int, obj interface{}) bool {
  if ctx.Config.SuccessFormat != nil {
    ctx.GinContext.JSON(code, ctx.Config.SuccessFormat(code, obj))
    ctx.GinContext.Abort()
    return true
  }
  ctx.GinContext.JSON(code, obj)
  ctx.GinContext.Abort()
  return true
}

func (ctx *Context) Success200(obj interface{}) bool {
  return ctx.Success(200, obj)
}

func (ctx *Context) Error(code int, errCode string, obj ...interface{}) bool {
  if ctx.Config.ErrorFormat != nil {
    errMessage := ""
    if ctx.Config.ErrorCodes != nil {
      errMessage = ctx.Config.ErrorCodes[errCode]
    }
    ctx.GinContext.JSON(code, ctx.Config.ErrorFormat(code, errCode, errMessage, obj))
    ctx.GinContext.Abort()
    return false
  }
  ctx.GinContext.JSON(code, obj)
  ctx.GinContext.Abort()
  return false
}

func (ctx *Context) Error200(errCode string, obj ...interface{}) bool {
  return ctx.Error(200, errCode, obj)
}

func (ctx *Context) Next() {
  ctx.GinContext.Next()
}

func (ctx *Context) Abourt() {
  ctx.GinContext.Abort()
}

func (ctx *Context) Set(key string, value interface{}) {
  ctx.GinContext.Set(key, value)
}

func (ctx *Context) Get(key string) (interface{}, bool) {
  return ctx.GinContext.Get(key)
}
