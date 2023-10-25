package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type (
	baseMiddleware struct {
		log logger.Logger
	}
	router interface {
		gin.IRouter

		GetStorage() models.Storage
		GetLogger() logger.Logger
	}
)

func Setup(r router) {
	bm := &baseMiddleware{
		log: r.GetLogger(),
	}

	r.Use(bm.Logger)
	r.Use(bm.Compress)
	r.Use(r.GetStorage().GetMiddleware())
}
