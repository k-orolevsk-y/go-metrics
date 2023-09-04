package handlers

import (
	"github.com/gin-gonic/gin"
)

func ValidateContentType(ctx *gin.Context, contentType string) bool {
	requestContentType := ctx.GetHeader("Content-Type")
	if requestContentType == "" {
		return true
	}

	return requestContentType == contentType
}
