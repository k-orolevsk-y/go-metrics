package router

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type Router struct {
	*gin.Engine

	storage models.Storage
	log     logger.Logger
}

func New(storage models.Storage, log logger.Logger) *Router {
	gin.SetMode(gin.ReleaseMode)

	return &Router{
		Engine: gin.New(),

		storage: storage,
		log:     log,
	}
}

func (r *Router) GetStorage() models.Storage {
	return r.storage
}

func (r *Router) GetLogger() logger.Logger {
	return r.log
}
