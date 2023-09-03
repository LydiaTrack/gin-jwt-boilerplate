package api

import (
	"gin-jwt-boilerplate/cmd/app/handlers"
	"github.com/gin-gonic/gin"
)

// InitHealth initializes health routes
func InitHealth(r *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	r.GET("/health", healthHandler.GetHealth)
}
