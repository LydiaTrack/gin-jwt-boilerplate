package api

import (
	"gin-jwt-boilerplate/cmd/app/handlers"
	"gin-jwt-boilerplate/internal/auth"
	"gin-jwt-boilerplate/internal/repository"
	"gin-jwt-boilerplate/internal/service"
	"github.com/gin-gonic/gin"
)

// InitAuth initializes auth routes
func InitAuth(r *gin.Engine) {
	userRepository := repository.GetUserRepository()
	userService := service.NewUserService(userRepository)
	sessionRepository := repository.GetSessionRepository()
	sessionService := service.NewSessionService(sessionRepository, userService)
	authService := auth.NewAuthService(userService, sessionService)

	authHandler := handlers.NewAuthHandler(authService)

	r.POST("/login", authHandler.Login)
	r.GET("/logout", authHandler.Logout)
	r.GET("/currentUser", authHandler.GetCurrentUser)
	r.POST("/refreshToken", authHandler.RefreshToken)
}
