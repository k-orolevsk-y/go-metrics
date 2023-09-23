package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (bh baseHandler) BadRequest(ctx *gin.Context) {
	bh.handleBadRequest(ctx)
}

func (bh baseHandler) handleBadRequest(ctx *gin.Context) {
	ctx.Status(http.StatusBadRequest)
	ctx.Abort()
}
