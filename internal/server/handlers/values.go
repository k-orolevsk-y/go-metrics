package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func (bh baseHandler) Values() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		values := bh.storage.GetAll()

		text := "<center><h1>Values</h1>"
		for _, value := range values {
			text += fmt.Sprintf("<p>%s: %s - %v</p>", value.Type, value.Name, value.Value)
		}
		text += "</center>"

		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")

		if _, err := io.Copy(ctx.Writer, strings.NewReader(text)); err != nil {
			bh.log.Errorf("io.Copy() error: %s", err)
			ctx.String(http.StatusInternalServerError, "%s", "Internal server error")
		}

		ctx.Abort()
	}
}
