package service

import (
	"net/http"
	"vizaduum/job"

	"github.com/gin-gonic/gin"
)

func Pause(ctx *gin.Context) {
	job.Status = 1 // 变更为维护状态
	ctx.JSON(http.StatusOK, gin.H{
		"message": "set pause",
	})
}

func Start(ctx *gin.Context) {
	job.Status = 0 // 解除维护状态
	ctx.JSON(http.StatusOK, gin.H{
		"message": "set start",
	})
}

func Status(ctx *gin.Context) {
	if job.Status == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "working",
			"status":  0,
		})
	} else if job.Status == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pause",
			"status":  1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "abnormal, please check",
			"status":  job.Status,
		})
	}
}
