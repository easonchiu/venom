package venom

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type IEngine interface {
	Engine() *gin.Engine
	Start() error
}

type Engine struct {
	c *Config
	g *gin.Engine
}

// 初始化server
func Init(config *Config) *Engine {
	setConfig(config)

	ginMode := gin.DebugMode
	if config.Mode == ProductionMode {
		ginMode = gin.ReleaseMode
	}

	gin.SetMode(ginMode)

	g := gin.Default()

	initRedisClient(DefaultRedisClientKey, config.Redis)
	initQmgoClient(DefaultQmgoClientKey, config.Qmgo)

	if config.RedisMap != nil {
		for k, c := range config.RedisMap {
			initRedisClient(k, c)
		}
	}

	if config.QmgoMap != nil {
		for k, c := range config.QmgoMap {
			initQmgoClient(k, c)
		}
	}

	if !config.Logger.Disabled {
		g.Use(LoggerMiddleware(config.Logger))
	}

	engine := &Engine{
		c: config,
		g: g,
	}

	return engine
}

func (e *Engine) Engine() *gin.Engine {
	return e.g
}

// 启动server
func (e *Engine) Start() error {
	defer func() {
		GetQmgoClient().CloseAll()
		GetRedisClient().CloseAll()
	}()

	fmt.Println("Ready start venom ...")
	return e.g.Run(e.c.Address + ":" + e.c.Port)
}

var (
	_ IEngine = &Engine{}
)
