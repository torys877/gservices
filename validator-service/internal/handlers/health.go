package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	db, err := h.db.DB()
	if err != nil {
		log.Println("Database connection error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}

	if err := db.Ping(); err != nil {
		log.Println("Database ping failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
