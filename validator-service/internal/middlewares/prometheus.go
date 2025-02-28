package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"validator-service/internal/monitoring"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()
		timer := prometheus.NewTimer(monitoring.ResponseDuration.WithLabelValues(endpoint))
		defer timer.ObserveDuration()

		monitoring.TotalRequests.WithLabelValues(endpoint).Inc()
		c.Next()
	}
}
