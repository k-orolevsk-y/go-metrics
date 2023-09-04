package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
	"net/http"
	"strconv"
)

func Value(storage stor.Storage) gin.HandlerFunc {
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

		var response interface{}
		storageType := ctx.Param("type")

		if storageType == string(stor.GaugeType) {
			value, err := storage.GetGauge(id)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			response = strconv.FormatFloat(value, 'f', -1, 64)
		} else if storageType == string(stor.CounterType) {
			value, err := storage.GetCounter(id)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			response = value
		} else {
			handleBadRequest(ctx)
			return
		}

		ctx.String(http.StatusOK, "%v", response)
		ctx.Abort()
	}
}
