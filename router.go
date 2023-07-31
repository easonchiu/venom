package venom

import "github.com/gin-gonic/gin"

type Router struct {
	URI         string // method + path，例： GET:/users/list
	Middlewares string // 中间件，多个中间件用英文逗号分隔，例：jwt,admin,xxx
	Handle      func(*gin.Context)
}
