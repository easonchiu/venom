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
  *gin.Context
}

func (ctx *Context) Bg() context.Context {
  return context.Background()
}

func (ctx *Context) Success(code int, obj interface{}) bool {
  if ctx.Config.SuccessFormat != nil {
    ctx.JSON(code, ctx.Config.SuccessFormat(code, obj))
    ctx.Abort()
    return true
  }
  ctx.JSON(code, obj)
  ctx.Abort()
  return true
}

func (ctx *Context) Success200(obj interface{}) bool {
  return ctx.Success(200, obj)
}

func (ctx *Context) Fail(code int, errCode string, obj ...interface{}) bool {
  if ctx.Config.FailFormat != nil {
    errMessage := ""
    if ctx.Config.FailCodes != nil {
      errMessage = ctx.Config.FailCodes[errCode]
    }
    ctx.JSON(code, ctx.Config.FailFormat(code, errCode, errMessage, obj))
    ctx.Abort()
    return false
  }
  ctx.JSON(code, obj)
  ctx.Abort()
  return false
}

func (ctx *Context) Fail200(errCode string, obj ...interface{}) bool {
  return ctx.Fail(200, errCode, obj)
}
