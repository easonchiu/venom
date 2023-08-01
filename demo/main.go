package main

import (
	"fmt"

	"github.com/easonchiu/venom"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	engine := venom.New(venom.DevelopmentMode, &venom.Config{
		Port: "3000",
	})

	engine.RegisterPlugins(
		venom.InitQmgoPlugin(&venom.QmgoConfig{
			URI:      "mongodb://localhost:27017",
			Database: "test",
		}),
	)

	engine.RegisterMiddlewares(
		venom.InitLoggerMiddleware(&venom.LoggerConfig{
			Name:     "global:log",
			Filename: "log",
			Level:    logrus.DebugLevel,
		}),
	)

	engine.RegisterRouters(
		venom.NewRouter("GET", "/list", testHandle),
		venom.NewRouterGroupWithMiddlewares("/group2", "auth22222",
			venom.NewRouterGroupWithMiddlewares("/group2", "auth3333",
				venom.NewRouterWithMiddlewares("GET", "/list", "auth4444", testHandle),
			),
		),
		venom.NewRouterWithMiddlewares("GET", "/list2", "auth", testHandle),
		venom.NewRouterGroup("/group1",
			venom.NewRouter("GET", "/list", testHandle),
			venom.NewRouterGroup("/group11",
				venom.NewRouter("GET", "/list", testHandle),
			),
		),
	)

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
	// venom.Success(ctx, "200")
}
