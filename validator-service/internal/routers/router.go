package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"validator-service/internal/handlers"
)

func SetupRoutes(r *gin.Engine, h *handlers.Handler) {
	// Validator endpoints
	r.POST("/validators", h.CreateValidator)
	r.GET("/validators/:request_id", h.CheckRequestStatus)

	// Health check endpoint
	r.GET("/health", h.HealthCheck)

	// Metrics endpoints
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
