package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Logger(log logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		method := ctx.Request.Method
		uri := ctx.Request.URL
		duration := time.Since(start)
		statusCode := ctx.Writer.Status()
		size := ctx.Writer.Size()

		log.Infof("%s Request - URI: \"%s\" - LeadTime: %v - StatusCode: %d (%s) - BodySize: %d", method, uri, duration, statusCode, http.StatusText(statusCode), size)
	}
}

type logger interface {
	Infof(template string, args ...interface{})
}
