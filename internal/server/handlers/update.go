package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"
	"net/http"
	"strconv"
)

func (bh baseHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ValidateContentType(ctx, "text/plain") {
			bh.handleBadRequest(ctx)
			return
		}

		id := ctx.Param("name")
		if id == "" {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()

			return
		}

		storageType := ctx.Param("type")
		if storageType == string(storage.GaugeType) {
			value, err := strconv.ParseFloat(ctx.Param("value"), 64)
			if err != nil {
				bh.handleBadRequest(ctx)
				return
			}

			bh.storage.SetGauge(id, value)
		} else if storageType == string(storage.CounterType) {
			value, err := strconv.ParseInt(ctx.Param("value"), 0, 64)
			if err != nil {
				bh.handleBadRequest(ctx)
				return
			}

			bh.storage.AddCounter(id, value)
		} else {
			bh.handleBadRequest(ctx)
			return
		}

		ctx.Status(http.StatusOK)
		ctx.Abort()
	}
}
