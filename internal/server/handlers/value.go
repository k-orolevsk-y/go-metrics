package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"
	"net/http"
	"strconv"
)

func (bh baseHandler) Value() gin.HandlerFunc {
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

		var response interface{}
		storageType := ctx.Param("type")

		if storageType == string(storage.GaugeType) {
			value, err := bh.storage.GetGauge(id)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			response = strconv.FormatFloat(value, 'f', -1, 64)
		} else if storageType == string(storage.CounterType) {
			value, err := bh.storage.GetCounter(id)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			response = value
		} else {
			bh.handleBadRequest(ctx)
			return
		}

		ctx.String(http.StatusOK, "%v", response)
		ctx.Abort()
	}
}
