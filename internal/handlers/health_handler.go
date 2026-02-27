package handlers

import (
	"finance-hub-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health handles GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	response.SuccessResponse(c, http.StatusOK, "Service is healthy", gin.H{
		"status": "ok",
	})
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(c *gin.Context) {
	// You can add database connectivity check here
	response.SuccessResponse(c, http.StatusOK, "Service is ready", gin.H{
		"status": "ready",
	})
}
