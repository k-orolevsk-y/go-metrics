package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"io"
	"net/http"
	"strings"
)

func (bh baseHandler) Values() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		values, err := bh.storage.GetAll()
		if err != nil {
			bh.log.Debugf("Error get all metrics: %s", err)

			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()

			return
		}

		text := "<center><h1>Values</h1>"
		for _, value := range values {
			if value.MType == string(models.GaugeType) {
				text += fmt.Sprintf("<p>%s: %s - %v</p>", value.MType, value.ID, *value.Value)
			} else if value.MType == string(models.CounterType) {
				text += fmt.Sprintf("<p>%s: %s - %d</p>", value.MType, value.ID, *value.Delta)
			}
		}
		text += "</center>"

		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")

		if _, err := io.Copy(ctx.Writer, strings.NewReader(text)); err != nil {
			bh.log.Debugf("io.Copy() error: %s", err)
			ctx.String(http.StatusInternalServerError, "%s", "Internal server error")
		}

		ctx.Abort()
	}
}
