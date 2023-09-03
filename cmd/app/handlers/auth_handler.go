package handlers

import (
	"gin-jwt-boilerplate/internal/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) AuthHandler {
	return AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login
// @Description login.
// @Tags auth
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func (h AuthHandler) Login(c *gin.Context) {
	var loginCommand auth.Request
	if err := c.ShouldBindJSON(&loginCommand); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Login(loginCommand)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description get current user.
// @Tags auth
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /currentUser [get]
func (h AuthHandler) GetCurrentUser(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description refresh token.
// @Tags auth
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /refreshToken [get]
func (h AuthHandler) RefreshToken(c *gin.Context) {
	tokenPair, err := h.authService.RefreshTokenPair(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tokenPair)
}

// Logout godoc
// @Summary Logout
// @Description logout.
// @Tags auth
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /logout [get]
func (h AuthHandler) Logout(c *gin.Context) {
	err := h.authService.Logout(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
