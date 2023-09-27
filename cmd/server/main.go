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

	defer func(logger *zap.Logger) {
		if err = logger.Sync(); err != nil {
			panic(err)
		}
	}(logger)

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

	r.POST("/value", baseHandler.ValueByBody())
	r.POST("/value/", baseHandler.ValueByBody())

	r.GET("/value/:type/:name", baseHandler.ValueByURI())
	r.GET("/value/:type/:name/", baseHandler.ValueByURI())

	r.POST("/update", baseHandler.UpdateByBody())
	r.POST("/update/", baseHandler.UpdateByBody())

	r.POST("/update/:type", baseHandler.UpdateByURI())
	r.POST("/update/:type/", baseHandler.UpdateByURI())

	r.POST("/update/:type/:name/:value", baseHandler.UpdateByURI())
	r.POST("/update/:type/:name/:value/", baseHandler.UpdateByURI())

	r.NoRoute(baseHandler.BadRequest)

	return r
}
