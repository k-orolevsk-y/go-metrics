package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (bh baseHandler) BadRequest(ctx *gin.Context) {
	bh.handleBadRequest(ctx)
}

func (bh baseHandler) handleBadRequest(ctx *gin.Context) {
	ctx.Status(http.StatusBadRequest)
	ctx.Abort()
}
