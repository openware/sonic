package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPublicConfigs returns public configs
func GetPublicConfigs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, memoryCache)
}
