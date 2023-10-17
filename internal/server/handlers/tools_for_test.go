package handlers

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/middlewares"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	serverRouter "github.com/k-orolevsk-y/go-metricts-tpl/internal/server/router"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

func setupRouter(storage models.Storage, log logger.Logger) *serverRouter.Router {
	r := serverRouter.New(storage, log)
	middlewares.Setup(r)
	Setup(r)

	return r
}

func getPointerFloat64(v float64) *float64 {
	return &v
}

func getPointerInt64(v int64) *int64 {
	return &v
}
