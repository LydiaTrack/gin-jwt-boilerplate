package service

import (
	"errors"
	"gin-jwt-boilerplate/internal/domain"
	"gin-jwt-boilerplate/internal/domain/commands"
	"gin-jwt-boilerplate/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

type UserRepository interface {
	// SaveUser saves a user
	SaveUser(user domain.UserModel) (domain.UserModel, error)
	// GetUser gets a user by id
	GetUser(id bson.ObjectId) (domain.UserModel, error)
	// GetUserByUsername gets a user by username
	GetUserByUsername(username string) (domain.UserModel, error)
	// ExistsUser checks if a user exists
	ExistsUser(id bson.ObjectId) (bool, error)
	// DeleteUser deletes a user by id
	DeleteUser(id bson.ObjectId) error
	// ExistsByUsername checks if a user exists by username
	ExistsByUsername(username string) bool
}

// CreateUser TODO: Add permission check
func (s UserService) CreateUser(command commands.CreateUserCommand) (domain.UserModel, error) {
	// TODO: These kind of operations must be done with specific requests, not by UserModel model itself
	// Validate user
	// Map command to user
	user := domain.NewUser(bson.NewObjectId().Hex(), command.Username,
		command.Password, command.PersonInfo, time.Now(), 1)
	if err := user.Validate(); err != nil {
		return user, err
	}
	userExists := s.userRepository.ExistsByUsername(user.Username)

	if userExists {
		return domain.UserModel{}, errors.New("user already exists")
	}

	user, err := beforeCreateUser(user)

	savedUser, err := s.userRepository.SaveUser(user)
	if err != nil {
		return domain.UserModel{}, err
	}
	savedUser, err = afterCreateUser(savedUser)
	if err != nil {
		return domain.UserModel{}, err
	}

	savedUser, _ = s.GetUser(savedUser.ID.Hex())
	//	utils.Log("User %s created successfully", savedUser.Username)
	return savedUser, nil
}

// beforeCreateUser is a hook that is called before creating a user
func beforeCreateUser(user domain.UserModel) (domain.UserModel, error) {
	// Hash user password before saving
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return domain.UserModel{}, err
	}

	user.Password = hashedPassword
	return user, nil
}

// afterCreateUser is a hook that is called after creating a user
func afterCreateUser(user domain.UserModel) (domain.UserModel, error) {
	return user, nil
}

func (s UserService) GetUser(id string) (domain.UserModel, error) {
	user, err := s.userRepository.GetUser(bson.ObjectIdHex(id))
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}

func (s UserService) ExistsUser(id string) (bool, error) {
	exists, err := s.userRepository.ExistsUser(bson.ObjectIdHex(id))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s UserService) DeleteUser(command commands.DeleteUserCommand) error {
	existsUser, err := s.ExistsUser(command.ID.Hex())
	if err != nil {
		return err
	}
	if !existsUser {
		return errors.New("user does not exist")
	}

	err = s.userRepository.DeleteUser(command.ID)
	if err != nil {
		return err
	}
	return nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyUser verifies a user by username and password
func (s UserService) VerifyUser(username string, password string) (domain.UserModel, error) {
	// Get the user by username
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return domain.UserModel{}, err
	}

	// Compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		utils.LogFatal("Error comparing passwords: " + err.Error())
		return domain.UserModel{}, err
	}

	return user, nil
}

// ExistsByUsername gets a user by username
func (s UserService) ExistsByUsername(username string) bool {
	return s.userRepository.ExistsByUsername(username)
}
