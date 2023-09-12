package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/handlers"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"
)

func main() {
	config.Load()
	if err := config.Parse(); err != nil {
		panic(err)
	}

	storage := stor.NewMem()

	r := setupRouter(&storage)
	if err := r.Run(config.Config.Address); err != nil {
		panic(err)
	}
}

func setupRouter(storage handlers.Storage) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	baseHandler := handlers.NewBase(storage)

	r.GET("/", baseHandler.Values())

	r.GET("/value/:type/:name", baseHandler.Value())
	r.GET("/value/:type/:name/", baseHandler.Value())

	// Gin не считает ссылки /update/gauge/ подходящими под условие /update/:type/:name/:value,
	// поэтому нужен такой костыль :(
	r.POST("/update/:type", baseHandler.Update())
	r.POST("/update/:type/", baseHandler.Update())

	r.POST("/update/:type/:name/:value", baseHandler.Update())
	r.POST("/update/:type/:name/:value/", baseHandler.Update())

	r.NoRoute(handlers.BadRequest)

	return r
}
