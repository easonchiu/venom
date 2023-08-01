package venom

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	config            *Config
	__gin             *gin.Engine
	__plugins         []IPlugin
	__middlewares     []IMiddleware
	__onBeforeStart   func() // 启动前的生命周期函数
	__onBeforeDestroy func() // 销毁前的生命周期函数
}

type IEngine interface {
	Start() error
	GinEngine() *gin.Engine
	BeforeStart(func())
	BeforeDestroy(func())
}

// 初始化server
func New(mode Mode, config *Config) *Engine {
	ginMode := gin.DebugMode
	if mode == ProductionMode {
		ginMode = gin.ReleaseMode
	}

	gin.SetMode(ginMode)

	g := gin.Default()

	engine := &Engine{
		config: config,
		__gin:  g,
	}

	return engine
}

// 初始化插件
func (engine *Engine) RegisterPlugins(plugins ...IPlugin) {
	engine.__plugins = plugins

	if plugins == nil {
		return
	}

	for _, plugin := range plugins {
		plugin.OnStart(engine.config)
	}
}

// 初始化中间件
func (engine *Engine) RegisterMiddlewares(middlewares ...IMiddleware) {
	engine.__middlewares = middlewares

	if middlewares == nil {
		return
	}

	for _, middleware := range middlewares {
		middleware.OnStart(engine.config)

		// 全局的中间件启用
		if strings.HasPrefix(strings.ToLower(middleware.Name()), "global:") {
			fmt.Printf("[VENOM] MIDDLEWARE - Use global middleware -> %v \n", middleware.Name())
			engine.__gin.Use(middleware.GetGinMiddleware(engine.config))
		}
	}
}

// 初始化路由
func (engine *Engine) RegisterRouters(routers ...Router) {
	engine.__registerRouters("", nil, routers...)
}

func (engine *Engine) __registerRouters(prefix string, middlewares []string, routers ...Router) {
	for _, router := range routers {
		mergedMiddlewares := middlewares
		if len(router.Middlewares) > 0 {
			mergedMiddlewares = append(mergedMiddlewares, router.Middlewares...)
		}

		if router.IsRouter() {
			mws := engine.GetMiddlewares(middlewares)
			handles := make([]gin.HandlerFunc, len(mergedMiddlewares)+1)
			for _, mw := range mws {
				handles = append(handles, mw.GetGinMiddleware(engine.config))
			}
			handles = append(handles, router.Handle)
			engine.GinEngine().Handle(router.Method, prefix+router.Path, handles...)
		} else if router.IsGroup() {
			engine.__registerRouters(prefix+router.Path, mergedMiddlewares, router.children...)
		}
	}
}

// 初始化server
// func Init(config *Config) *Engine {
// 	SetConfig(config)

// 	ginMode := gin.DebugMode
// 	if config.Mode == ProductionMode {
// 		ginMode = gin.ReleaseMode
// 	}

// 	gin.SetMode(ginMode)

// 	g := gin.Default()

// 	// use 中间件
// 	// 当中间件以 global: 开头时，全局引用
// 	for name, middleware := range config.Middlewares {
// 		if strings.HasPrefix(strings.ToLower(name), "global:") {
// 			g.Use(middleware.GetGinMiddleware(config))
// 		}
// 	}

// 	// 注册路由
// 	for _, route := range config.Routers {
// 		spuri := strings.SplitN(route.URI, ":", 2)
// 		if len(spuri) != 2 {
// 			break
// 		}

// 		handlers := make([]gin.HandlerFunc, 0)
// 		spmw := strings.Split(route.Middlewares, ",")

// 		// 在 MiddlewarePrefix 中有匹配的路由时，在路由 slice 最前方加入该中间件
// 		for prefix, mw := range config.MiddlewarePrefix {
// 			if strings.HasPrefix(spuri[1], prefix) {
// 				spmw = append([]string{mw}, spmw...)
// 			}
// 		}

// 		// 循环中间件，找到具体的中间件方法
// 		for _, mw := range spmw {
// 			if mwFun, ok := config.Middlewares[mw]; ok {
// 				handlers = append(handlers, mwFun.GetGinMiddleware(config))
// 			}
// 		}

// 		// 在 handlers 最后加处 handle
// 		handlers = append(handlers, route.Handle)

// 		// 在 gin 中加入路由
// 		g.Handle(strings.ToUpper(spuri[0]), spuri[1], handlers...)
// 	}

// 	engine := &Engine{
// 		c: config,
// 		g: g,
// 	}

// 	return engine
// }

func (e *Engine) GinEngine() *gin.Engine {
	return e.__gin
}

// 启动server
func (e *Engine) Start() error {
	defer func(config *Config) {
		// 调用生命周期函数
		if e.__onBeforeDestroy != nil {
			e.__onBeforeDestroy()
		}

		// 销毁中间件
		for _, middleware := range e.__middlewares {
			middleware.OnDestroy(config)
		}

		// 销毁插件
		for _, plugin := range e.__plugins {
			plugin.OnDestroy(config)
		}
	}(e.config)

	// 调用生命周期函数
	if e.__onBeforeStart != nil {
		e.__onBeforeStart()
	}

	fmt.Println("[VENOM] 🎉🎉🎉 Ready start venom ...")

	return e.__gin.Run(e.config.Address + ":" + e.config.Port)
}

// 保存生命周期函数
func (e *Engine) BeforeStart(f func()) {
	e.__onBeforeStart = f
}

// 保存生命周期函数
func (e *Engine) BeforeDestroy(f func()) {
	e.__onBeforeDestroy = f
}

var (
	_ IEngine = &Engine{}
)
