package main

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "venom"
  "venom/example/controller"
)

func main() {
  v := venom.Init(&venom.Config{
    Mode: venom.ModeDevelopment,
    Port: "3200",
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
      "999999": "system error",
    },
    Redis: venom.RedisConfig{
      Host:     "localhost",
      Port:     "6379",
      Password: "",
    },
    Mongo: venom.MongoConfig{
      URI: "mongodb://localhost:27017/db",
    },
    Logger: venom.LoggerConfig{
      Filename:   "system.log",
      MaxSize:    500,
      MaxBackups: 3,
      MaxAge:     30,
      Level:      logrus.DebugLevel,
    },
  })

  RegisterRouters(v.Router())

  _ = v.Start()
}

func RegisterRouters(r *venom.Router) {
  user := new(controller.UserController)
  r.GET("/", user.Get)
}
