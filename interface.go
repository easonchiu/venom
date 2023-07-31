package venom

import "github.com/gin-gonic/gin"

type IEngine interface {
	Start() error
	GinEngine() *gin.Engine
	BeforeStart(func())
	BeforeDestroy(func())
}

type IPlugin interface {
	OnStart(*Config)
	OnDestroy(*Config)
}

type IMiddleware interface {
	OnStart(*Config)
	OnDestroy(*Config)
	GetGinMiddleware(*Config) gin.HandlerFunc
}
