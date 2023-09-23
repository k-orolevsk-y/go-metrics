package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/handlers"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/middlewares"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugarLogger := logger.Sugar()

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Panic loading config: %s", err)
	}

	memStorage := storage.NewMem()

	r := setupRouter(&memStorage, sugarLogger)
	if err = r.Run(config.Config.Address); err != nil {
		sugarLogger.Panicf("Panic start server: %s", err)
	}
}

func setupRouter(storage *storage.Mem, logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middlewares.Logger(logger))

	baseHandler := handlers.NewBase(storage, logger)

	r.GET("/", baseHandler.Values())

	r.GET("/value/:type/:name", baseHandler.Value())
	r.GET("/value/:type/:name/", baseHandler.Value())

	// Gin не считает ссылки /update/gauge/ подходящими под условие /update/:type/:name/:value,
	// поэтому нужен такой костыль :(
	r.POST("/update/:type", baseHandler.Update())
	r.POST("/update/:type/", baseHandler.Update())

	r.POST("/update/:type/:name/:value", baseHandler.Update())
	r.POST("/update/:type/:name/:value/", baseHandler.Update())

	r.NoRoute(baseHandler.BadRequest)

	return r
}
