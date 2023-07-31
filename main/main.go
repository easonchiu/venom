package main

import (
	"fmt"

	"github.com/easonchiu/venom"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	engine := venom.Init(&venom.Config{
		Port:    "5000",
		Mode:    venom.DevelopmentMode,
		Plugins: []venom.IPlugin{},
		Middlewares: []venom.IMiddleware{
			venom.InitLoggerMiddleware(&venom.LoggerConfig{
				Filename: "log",
				Level:    logrus.DebugLevel,
			}),
		},
	})

	engine.BeforeStart(func() {
		fmt.Println("before start...")
	})

	engine.BeforeDestroy(func() {
		fmt.Println("before destroy...")
	})

	g := engine.GinEngine()

	g.GET("/", testHandle)

	if err := engine.Start(); err != nil {
		panic(err)
	}
}

func testHandle(ctx *gin.Context) {
	venom.Success(ctx, "200")
}
