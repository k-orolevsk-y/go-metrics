package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type (
	baseHandler struct {
		storage models.Storage
		log     logger.Logger
	}
	router interface {
		gin.IRouter
		NoRoute(...gin.HandlerFunc)

		GetStorage() models.Storage
		GetLogger() logger.Logger
	}
)

func Setup(r router) {
	bh := &baseHandler{storage: r.GetStorage(), log: r.GetLogger()}

	r.GET("/", bh.Values())

	r.GET("/ping", bh.Ping())

	r.POST("/value", bh.ValueByBody())
	r.POST("/value/", bh.ValueByBody())

	r.GET("/value/:type/:name", bh.ValueByURI())
	r.GET("/value/:type/:name/", bh.ValueByURI())

	r.POST("/updates", bh.Updates())
	r.POST("/updates/", bh.Updates())

	r.POST("/update", bh.UpdateByBody())
	r.POST("/update/", bh.UpdateByBody())

	r.POST("/update/:type", bh.UpdateByURI())
	r.POST("/update/:type/", bh.UpdateByURI())
	r.POST("/update/:type/:name/:value", bh.UpdateByURI())
	r.POST("/update/:type/:name/:value/", bh.UpdateByURI())

	r.NoRoute(bh.BadRequest)
}
