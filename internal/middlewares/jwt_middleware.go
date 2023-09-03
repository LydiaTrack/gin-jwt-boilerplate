package middlewares

import (
	"gin-jwt-boilerplate/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware is a middleware for JWT authentication
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.ExtractTokenFromContext(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		err = auth.IsTokenValid(token)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
