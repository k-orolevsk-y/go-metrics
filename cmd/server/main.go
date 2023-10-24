package main

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/handlers"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/middlewares"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/router"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

func main() {
	sugarLogger, err := logger.New()
	if err != nil {
		panic(err)
	}
	sugarLogger.Debugf("The logger has been successfully initialized and configured.")

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Failed loading config: %s", err)
	}
	sugarLogger.Debugf("The config was successfully received and configured.")

	store, err := storage.Setup(sugarLogger)
	if err != nil {
		sugarLogger.Panicf("Failed setup storage: %s", err)
	}
	sugarLogger.Debugf("Selected storage: %s", store)

	defer func() {
		if err = sugarLogger.Sync(); err != nil {
			panic(err)
		}

		if err = store.Close(); err != nil {
			panic(err)
		}
	}()

	r := router.New(store, sugarLogger)
	middlewares.Setup(r)
	handlers.Setup(r)

	sugarLogger.Debugf("Server routing is configured and sent to launch on: %s", config.Config.Address)
	if err = r.Run(config.Config.Address); err != nil {
		sugarLogger.Panicf("Failed start server: %s", err)
	}
}
