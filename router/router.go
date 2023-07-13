package router

import (
	"vizaduum/service"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	apiv1 := r.Group("/api/v1")
	apiv1.GET("/display", service.Display)

	return r
}
