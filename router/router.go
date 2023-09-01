package router

import (
	"vizaduum/service"

	"github.com/gin-gonic/gin"
)

// Setup initializes and configures a new gin.Engine instance.
//
// It sets up the necessary middlewares for logging and recovery.
// It also defines the routes for the "/api/v1" group, including the
// "/display", "/pause", "/start", and "/status" endpoints.
//
// Returns a pointer to the gin.Engine instance.
func Setup() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	apiv1 := r.Group("/api/v1")
	apiv1.GET("/health", service.HealthCheck)
	apiv1.GET("/display", service.Display)
	apiv1.POST("/pause", service.Pause)
	apiv1.POST("/start", service.Start)
	apiv1.GET("/status", service.Status)
	return r
}
