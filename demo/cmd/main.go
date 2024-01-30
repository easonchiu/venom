package main

import (
	"github.com/easonchiu/venom"
	"github.com/easonchiu/venom/demo/internal/handler"
	"github.com/easonchiu/venom/demo/internal/repo"
	"github.com/easonchiu/venom/demo/internal/router"
	"github.com/easonchiu/venom/demo/internal/service"
	"github.com/easonchiu/venom/demo/pkg/config"
)

func main() {
	conf := config.Load("...")

	engine := venom.New(&venom.Config{
		Mode: conf.Mode,
		Port: conf.Port,
	})

	engine.RegisterPlugins(
		venom.InitQmgoPlugin(&venom.QmgoConfig{
			URI:      "mongodb://localhost:27017",
			Database: "test",
		}),
	)

	initRouter(engine)

	engine.Start()
}

func initRouter(engine *venom.Engine) {
	db := venom.GetQmgoDB()
	userRepo := repo.NewUserRepo(db.Collection("user"))
	userService := service.NewUserService(userRepo)
	userHandle := handler.NewUserHandler(userService)
	router.NewUserRouter(userHandle).Load(engine)
}
