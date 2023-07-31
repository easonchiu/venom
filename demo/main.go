package main

import (
	"fmt"

	"github.com/easonchiu/venom"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	engine := venom.Init(&venom.Config{
		Port:    "3000",
		Mode:    venom.DevelopmentMode,
		Plugins: []venom.IPlugin{},
		Middlewares: map[string]venom.IMiddleware{
			"global:log": venom.InitLoggerMiddleware(&venom.LoggerConfig{
				Filename: "log",
				Level:    logrus.DebugLevel,
			}),
		},
		MiddlewarePrefix: map[string]string{
			"/console": "auth",
		},
		Routers: []venom.Router{
			{URI: "GET:/console/a", Handle: testHandle},
		},
	})

	engine.BeforeStart(func() {
		fmt.Println("before start...")
	})

	engine.BeforeDestroy(func() {
		fmt.Println("before destroy...")
	})

	if err := engine.Start(); err != nil {
		panic(err)
	}
}

func testHandle(ctx *gin.Context) {
	venom.Success(ctx, "200")
}
