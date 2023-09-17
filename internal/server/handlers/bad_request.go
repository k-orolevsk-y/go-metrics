package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BadRequest(ctx *gin.Context) {
	handleBadRequest(ctx)
}

func handleBadRequest(ctx *gin.Context) {
	ctx.Status(http.StatusBadRequest)
	ctx.Abort()
}
