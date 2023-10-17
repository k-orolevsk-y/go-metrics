package middlewares

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/handlers"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	serverRouter "github.com/k-orolevsk-y/go-metricts-tpl/internal/server/router"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

func setupRouter(storage models.Storage, log logger.Logger) *serverRouter.Router {
	r := serverRouter.New(storage, log)
	Setup(r)
	handlers.Setup(r)

	return r
}
