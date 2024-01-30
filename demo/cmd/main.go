package main

import (
	"github.com/easonchiu/venom"
	"github.com/easonchiu/venom/demo/internal/router"
	"github.com/easonchiu/venom/demo/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.Load("...")

	engine := venom.New(&venom.Config{
		Mode: conf.Mode,
		Port: conf.Port,
	})

	engine.RegisterPlugins(
		venom.InitQmgoPlugin(&venom.QmgoConfig{
			Key:      "test",
			URI:      "mongodb://localhost:27017",
			Database: "test",
		}),
		venom.InitQmgoPlugin(&venom.QmgoConfig{
			Key:      "test2",
			URI:      "mongodb://localhost:27017",
			Database: "test2",
		}),
	)

	router.NewMyRouter().Load(engine)

	engine.Start()
}

func testHandle(ctx *gin.Context) {
	// venom.Success(ctx, "200")
}
