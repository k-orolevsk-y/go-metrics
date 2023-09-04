package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/handlers"
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
)

func main() {
	storage := stor.NewMem()

	r := setupRouter(&storage)
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func setupRouter(storage stor.Storage) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/", handlers.Values(storage))

	r.GET("/value/:type/:name", handlers.Value(storage))
	r.GET("/value/:type/:name/", handlers.Value(storage))

	// Gin не считает ссылки /update/gauge/ подходящими под условие /update/:type/:name/:value,
	// поэтому нужен такой костыль :(
	r.POST("/update/:type", handlers.Update(storage))
	r.POST("/update/:type/", handlers.Update(storage))

	r.POST("/update/:type/:name/:value", handlers.Update(storage))
	r.POST("/update/:type/:name/:value/", handlers.Update(storage))

	r.NoRoute(handlers.BadRequest)

	return r
}
