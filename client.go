package venom

import "github.com/gin-gonic/gin"

type Client struct {
  Redis   string
  Mongodb string
  Gin     *gin.Engine
}
