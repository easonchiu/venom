/*
 * @Author: zhaozhida zhaozhida@qiniu.com
 * @Date: 2023-07-26 10:32:09
 * @LastEditors: zhaozhida zhaozhida@qiniu.com
 * @LastEditTime: 2023-07-26 14:25:11
 * @Description:
 */
package venom

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	c *Config
	g *gin.Engine
}

// 启动前的生命周期函数
var onBeforeStart func()

// 销毁前的生命周期函数
var onBeforeDestroy func()

// 初始化server
func Init(config *Config) *Engine {
	SetConfig(config)

	ginMode := gin.DebugMode
	if config.Mode == ProductionMode {
		ginMode = gin.ReleaseMode
	}

	gin.SetMode(ginMode)

	g := gin.Default()

	// use 中间件
	for _, middleware := range config.Middlewares {
		g.Use(middleware.GetGinMiddleware(config))
	}

	engine := &Engine{
		c: config,
		g: g,
	}

	return engine
}

func (e *Engine) GinEngine() *gin.Engine {
	return e.g
}

// 启动server
func (e *Engine) Start() error {
	defer func() {
		// 调用生命周期函数
		if onBeforeDestroy != nil {
			onBeforeDestroy()
		}

		// 销毁插件
		for _, plugin := range config.Plugins {
			plugin.OnDestroy(config)
		}
	}()

	// 启动插件
	for _, plugin := range config.Plugins {
		plugin.OnStart(config)
	}

	// 启动中间件
	for _, middleware := range config.Middlewares {
		middleware.OnStart(config)
	}

	// 调用生命周期函数
	if onBeforeStart != nil {
		onBeforeStart()
	}

	fmt.Println("🎉 Ready start venom ...")

	return e.g.Run(e.c.Address + ":" + e.c.Port)
}

// 保存生命周期函数
func (e *Engine) BeforeStart(f func()) {
	onBeforeStart = f
}

// 保存生命周期函数
func (e *Engine) BeforeDestroy(f func()) {
	onBeforeDestroy = f
}

var (
	_ IEngine = &Engine{}
)
