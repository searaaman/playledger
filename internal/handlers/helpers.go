package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseID reads a numeric path parameter and writes a 400 response if it is invalid.
func parseID(ctx *gin.Context, name string) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param(name), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid " + name,
		})
		return 0, false
	}
	return uint(id), true
}
