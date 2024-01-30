package venom

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
	config    *Config
	__plugins []IPlugin
}

type IEngine interface {
	Start()
}

var _ IEngine = (*Engine)(nil)

// New 初始化server
func New(config *Config) *Engine {
	gin.SetMode(config.Mode)
	return &Engine{gin.Default(), config, nil}
}

// RegisterPlugins 初始化插件
func (engine *Engine) RegisterPlugins(plugins ...IPlugin) {
	engine.__plugins = plugins

	if plugins == nil {
		return
	}

	for _, plugin := range plugins {
		plugin.OnStart(engine.config)
	}
}

// Start 启动server
func (e *Engine) Start() {
	defer func(config *Config) {
		// 销毁插件
		for _, plugin := range e.__plugins {
			plugin.OnDestroy(config)
		}
	}(e.config)

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", e.config.Address, e.config.Port),
		Handler: e,
	}

	// 启动 HTTP 服务器
	go func() {
		fmt.Printf("\n[VENOM] Server started at [%v]...\n\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// 创建一个信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待接收到信号
	<-quit
	log.Println("[VENOM] Shutting down server...")

	// 创建一个上下文，超时时间为 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("[VENOM] Server forced to shutdown:", err)
	}

	log.Println("[VENOM] Server stopped gracefully")
}
