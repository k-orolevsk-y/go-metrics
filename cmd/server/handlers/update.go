package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
	"net/http"
	"strconv"
)

func Update(storage stor.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ValidateContentType(ctx, "text/plain") {
			handleBadRequest(ctx)
			return
		}

		id := ctx.Param("name")
		if id == "" {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()

			return
		}

		storageType := ctx.Param("type")
		if storageType == string(stor.GaugeType) {
			value, err := strconv.ParseFloat(ctx.Param("value"), 64)
			if err != nil {
				handleBadRequest(ctx)
				return
			}

			storage.SetGauge(id, value)
		} else if storageType == string(stor.CounterType) {
			value, err := strconv.ParseInt(ctx.Param("value"), 0, 64)
			if err != nil {
				handleBadRequest(ctx)
				return
			}

			storage.AddCounter(id, value)
		} else {
			handleBadRequest(ctx)
			return
		}

		ctx.Status(http.StatusOK)
		ctx.Abort()
	}
}
