package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
)

func (bh baseHandler) Updates() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !bh.validateContentType(ctx, "application/json", false) {
			bh.log.Debugf("Request with invalid content-type.")
			bh.handleBadRequest(ctx)
			return
		}

		var objects []models.MetricsUpdate
		if response, statusCode, err := bh.validateAndShouldBindJSON(ctx, &objects); err != nil {
			if statusCode == http.StatusInternalServerError {
				bh.log.Errorf("Error decoding object request: %s (%T)", err, err)
			}

			if response == nil {
				ctx.Status(statusCode)
			} else {
				ctx.JSON(statusCode, response)
			}

			ctx.Abort()

			return
		}

		tx, err := bh.storage.NewTx()
		if err != nil {
			bh.log.Debugf("Failed to create transaction: %s (%T)", err, err)

			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()

			return
		}

		for _, obj := range objects {
			if obj.MType == string(models.GaugeType) {
				if err = tx.SetGauge(obj.ID, obj.Value); err != nil {
					bh.log.Errorf("Error set gauge (tx): %s (%T)", err, err)
					if err = tx.RollBack(); err != nil {
						bh.log.Errorf("Failed to rollback transaction [gauge]: %s (%T)", err, err)
					}

					ctx.Status(http.StatusInternalServerError)
					ctx.Abort()

					return
				}
			} else if obj.MType == string(models.CounterType) {
				if err = tx.AddCounter(obj.ID, obj.Delta); err != nil {
					bh.log.Errorf("Error add counter (tx): %s (%T)", err, err)
					if err = tx.RollBack(); err != nil {
						bh.log.Errorf("Failed to rollback transaction [counter]: %s (%T)", err, err)
					}

					ctx.Status(http.StatusInternalServerError)
					ctx.Abort()

					return
				}
			}
		}

		if err = tx.Commit(); err != nil {
			bh.log.Errorf("Failed to save changes from transaction: %s (%T)", err, err)

			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()

			return
		}

		ctx.Status(http.StatusOK)
		ctx.Abort()
	}
}
