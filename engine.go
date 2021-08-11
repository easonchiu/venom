package venom

import (
  "github.com/gin-gonic/gin"
)

type Engine struct {
  Config *Config
  Client *Client
  Gin    *gin.Engine
}

func Init(config *Config) *Engine {
  g := gin.Default()

  engine := &Engine{
    Config: config,
    Client: &Client{},
    Gin:    g,
  }

  return engine
}

func (e *Engine) Start() error {
  return e.Gin.Run(e.Config.Address + ":" + e.Config.Port)
}
