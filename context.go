package venom

import (
  "context"
  "errors"
  "fmt"
  "github.com/gin-gonic/gin"
)

type IContext interface {
  Bg() context.Context
  Success(code int, obj interface{}) bool
  Success200(obj interface{}) bool
  Fail(code int, errCode string, obj ...interface{}) bool
  Fail200(errCode string, obj ...interface{}) bool
}

type Context struct {
  Config *Config
  Redis  *RedisClient
  Mongo  *MongoClient
  Qmgo   *QmgoClient
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
  var data interface{} = nil
  if obj != nil && len(obj) > 0 {
    data = obj[0]
  }

  if ctx.Config.FailFormat != nil {
    errMessage := ""
    if ctx.Config.FailCodes != nil {
      errMessage = ctx.Config.FailCodes[errCode]
    }
    err := ctx.Config.FailFormat(code, errCode, errMessage, data)
    ctx.JSON(code, err)
    _ = ctx.Error(errors.New(fmt.Sprintf("status: %v, err_code: %v, data: %v", code, errCode, err)))
    ctx.Abort()
    return false
  }

  ctx.JSON(code, data)
  _ = ctx.Error(errors.New(fmt.Sprintf("status: %v, err_code: %v, data: %v", code, errCode, data)))
  ctx.Abort()
  return false
}

func (ctx *Context) Fail200(errCode string, obj ...interface{}) bool {
  return ctx.Fail(200, errCode, obj...)
}

var (
  _ IContext = &Context{}
)
