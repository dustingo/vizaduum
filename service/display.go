package service

import (
	"net/http"
	"vizaduum/config"

	"github.com/gin-gonic/gin"
)

// Display displays the JSON response for the given gin Context.
func Display(ctx *gin.Context) {
	data := config.GConfig.GameConfig
	ctx.JSON(http.StatusOK, gin.H{
		"messgae": "ok",
		"data": gin.H{
			"gameConfig": data,
		},
	})
}
