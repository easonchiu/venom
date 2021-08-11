package main

import (
  "github.com/gin-gonic/gin"
  "venom"
)

func main() {
  g := venom.Init(&venom.Config{
    Mode: venom.ModeDevelopment,
    Port: "3333",
    SuccessFormat: func(code int, obj interface{}) interface{} {
      return gin.H{
        "code":    "0",
        "message": "success",
        "data":    obj,
      }
    },
    ErrorFormat: func(code int, errCode string, errMessage string, obj ...interface{}) interface{} {
      return gin.H{
        "code":    errCode,
        "message": errMessage,
      }
    },
    ErrorCodes: map[string]string{
      "600123": "some error",
    },
    Redis: venom.RedisConfig{
      Address: "172.16.22.212",
      Port:    "6379",
    },
  })

  RegisterRouters(g.Router())

  _ = g.Start()
}

func RegisterRouters(r *venom.Router) {
  r.GET("/", Get2)
  r.Group("a").GET("/", Get)
}

func Get(ctx *venom.Ctx) bool {
  return ctx.Success200("res")
}

func Get2(ctx *venom.Ctx) bool {
  return ctx.Error200("600123")
}
