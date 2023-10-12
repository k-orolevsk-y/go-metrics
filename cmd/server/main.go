package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/database_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/file_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/handlers"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/middlewares"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/database"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

func main() {
	sugarLogger, err := logger.New()
	if err != nil {
		panic(err)
	}

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Failed loading config: %s", err)
	}

	storage, err := setupStorage(sugarLogger)
	if err != nil {
		sugarLogger.Panicf("Failed setup storage: %s", err)
	}

	defer func() {
		if err = sugarLogger.Sync(); err != nil {
			panic(err)
		}

		if err = storage.Close(); err != nil {
			panic(err)
		}
	}()

	r := setupRouter(storage, sugarLogger)
	if err = r.Run(config.Config.Address); err != nil {
		sugarLogger.Panicf("Failed start server: %s", err)
	}
}

func setupStorage(log logger.Logger) (models.Storage, error) {
	if config.Config.DatabaseDSN != "" {
		db, err := database.New()
		if err != nil {
			return nil, err
		}

		return dbstorage.New(db, log)
	} else if config.Config.FileStoragePath != "" {
		fs, err := filestorage.New(log)
		if err != nil {
			return nil, err
		}

		if err = fs.Restore(); err != nil {
			return nil, err
		}
		fs.Start()

		return fs, nil
	} else {
		return memstorage.NewMem(), nil
	}
}

func setupRouter(storage models.Storage, log logger.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	baseHandler := handlers.NewBase(storage, log)
	baseMiddleware := middlewares.NewBase(log)

	r.Use(baseMiddleware.Logger)
	r.Use(baseMiddleware.Compress)
	r.Use(storage.GetMiddleware())

	r.GET("/", baseHandler.Values())

	r.GET("/ping", baseHandler.Ping())

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
