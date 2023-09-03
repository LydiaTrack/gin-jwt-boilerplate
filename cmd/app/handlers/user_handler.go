package handlers

import (
	"gin-jwt-boilerplate/internal/domain/commands"
	"gin-jwt-boilerplate/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService: userService}
}

// GetUser godoc
// @Summary Get user by ID
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users:id [get]
func (h UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

// CreateUser godoc
// @Summary Create user
// @Description create user.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users [post]
func (h UserHandler) CreateUser(c *gin.Context) {
	var createUserCommand commands.CreateUserCommand
	if err := c.ShouldBindJSON(&createUserCommand); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userService.CreateUser(createUserCommand)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description delete user.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users:id [delete]
func (h UserHandler) DeleteUser(c *gin.Context) {
	var deleteUserCommand commands.DeleteUserCommand
	if err := c.ShouldBindJSON(&deleteUserCommand); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := h.userService.DeleteUser(deleteUserCommand)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
