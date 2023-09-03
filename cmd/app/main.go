package main

import (
	"gin-jwt-boilerplate/cmd/app/api"
	"gin-jwt-boilerplate/internal/service"
	"gin-jwt-boilerplate/internal/utils"

	"github.com/gin-gonic/gin"
)

// @title Main API
// @version 0.0.1
// @description Main API

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	r := gin.New()

	// Initialize routes
	initializeRoutes(r)
	// Initialize logging
	utils.InitLogging()
	// Initialize default user
	service.InitializeDefaultUser()

	// Run server on port 8080
	r.Run(":8080")
}

// initializeRoutes initializes routes for each API
func initializeRoutes(r *gin.Engine) {
	// If you want to add global interceptors, add them to this slice
	globalInterceptors := []gin.HandlerFunc{gin.Recovery(), gin.Logger()}

	r.Use(globalInterceptors...)

	api.InitUser(r)
	api.InitAuth(r)
	api.InitHealth(r)
}
