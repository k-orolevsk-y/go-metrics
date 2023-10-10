package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"net/http"
	"time"
)

func (bh baseHandler) Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxDB, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := bh.db.PingContext(ctxDB); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: fmt.Sprint(err),
			})
		} else {
			ctx.Status(http.StatusOK)
		}

		ctx.Abort()
	}
}
