package api

import (
	"gin-jwt-boilerplate/cmd/app/handlers"
	"gin-jwt-boilerplate/internal/middlewares"
	"gin-jwt-boilerplate/internal/repository"
	"gin-jwt-boilerplate/internal/service"
	"github.com/gin-gonic/gin"
)

// InitUser initializes user routes
func InitUser(r *gin.Engine) {

	userRepository := repository.GetUserRepository()
	userService := service.NewUserService(userRepository)

	userHandler := handlers.NewUserHandler(userService)

	authorizedPath := r.Group("/users")
	authorizedPath.Use(middlewares.JwtAuthMiddleware()).
		POST("", userHandler.CreateUser).
		GET("/:id", userHandler.GetUser).
		DELETE("/:id", userHandler.DeleteUser)
}
