package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	stor "github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
	"io"
	"net/http"
	"strings"
)

func Values(storage stor.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		values := storage.GetAll()

		text := "<center><h1>Values</h1>"
		for _, value := range values {
			text += fmt.Sprintf("<p>%s: %s - %v</p>", value.Type, value.Name, value.Value)
		}
		text += "</center>"

		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")

		if _, err := io.Copy(ctx.Writer, strings.NewReader(text)); err != nil {
			ctx.String(http.StatusInternalServerError, "%s", "Internal server error")
		}

		ctx.Abort()
	}
}
