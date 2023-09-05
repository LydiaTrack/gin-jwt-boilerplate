package service

import (
	"gin-jwt-boilerplate/internal/domain"
	"gin-jwt-boilerplate/internal/domain/commands"
	"gin-jwt-boilerplate/internal/repository"
	"gin-jwt-boilerplate/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

// InitializeDefaultUser initializes the default user with default credentials
func InitializeDefaultUser() error {
	userCreateCmd := commands.CreateUserCommand{
		Username: os.Getenv("DEFAULT_USER_USERNAME"),
		Password: os.Getenv("DEFAULT_USER_PASSWORD"),
		PersonInfo: domain.PersonInfo{
			FirstName: "Default",
			LastName:  "User",
			BirthDate: primitive.
				DateTime(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)),
		},
	}

	_, err := NewUserService(repository.GetUserRepository()).
		CreateUser(userCreateCmd)
	if err != nil {
		return err
	}

	utils.Log("Default user created successfully")

	return nil
}
