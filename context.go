package venom

import "github.com/gin-gonic/gin"

type Ctx struct {
  Config     *Config
  Client     *Client
  GinContext *gin.Context
}

func (ctx *Ctx) Success(code int, obj interface{}) bool {
  if ctx.Config.SuccessFormat != nil {
    ctx.GinContext.JSON(code, ctx.Config.SuccessFormat(code, obj))
    return true
  }
  ctx.GinContext.JSON(code, obj)
  return true
}

func (ctx *Ctx) Success200(obj interface{}) bool {
  return ctx.Success(200, obj)
}

func (ctx *Ctx) Error(code int, errCode string, obj... interface{}) bool {
  if ctx.Config.ErrorFormat != nil {
    errMessage := ""
    if ctx.Config.ErrorCodes != nil {
      errMessage = ctx.Config.ErrorCodes[errCode]
    }
    ctx.GinContext.JSON(code, ctx.Config.ErrorFormat(code, errCode, errMessage, obj))
    return false
  }
  ctx.GinContext.JSON(code, obj)
  return false
}

func (ctx *Ctx) Error200(errCode string, obj... interface{}) bool {
  return ctx.Error(200, errCode, obj)
}
