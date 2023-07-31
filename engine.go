package venom

import (
	"fmt"
	"strings"

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
	// 当中间件以 global: 开头时，全局引用
	for name, middleware := range config.Middlewares {
		if strings.HasPrefix(strings.ToLower(name), "global:") {
			g.Use(middleware.GetGinMiddleware(config))
		}
	}

	// 注册路由
	for _, route := range config.Routers {
		spuri := strings.SplitN(route.URI, ":", 2)
		if len(spuri) != 2 {
			break
		}

		handlers := make([]gin.HandlerFunc, 0)
		spmw := strings.Split(route.Middlewares, ",")

		// 在 MiddlewarePrefix 中有匹配的路由时，在路由 slice 最前方加入该中间件
		for prefix, mw := range config.MiddlewarePrefix {
			if strings.HasPrefix(spuri[1], prefix) {
				spmw = append([]string{mw}, spmw...)
			}
		}

		// 循环中间件，找到具体的中间件方法
		for _, mw := range spmw {
			if mwFun, ok := config.Middlewares[mw]; ok {
				handlers = append(handlers, mwFun.GetGinMiddleware(config))
			}
		}

		// 在 handlers 最后加处 handle
		handlers = append(handlers, route.Handle)

		// 在 gin 中加入路由
		g.Handle(strings.ToUpper(spuri[0]), spuri[1], handlers...)
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
