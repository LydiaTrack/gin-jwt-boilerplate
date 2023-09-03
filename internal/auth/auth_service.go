package auth

import (
	"errors"
	"gin-jwt-boilerplate/internal/domain"
	"gin-jwt-boilerplate/internal/domain/commands"
	"gin-jwt-boilerplate/internal/service"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

type UserService interface {
	ExistsByUsername(username string) bool
	VerifyUser(username string, password string) (domain.UserModel, error)
	GetUser(id string) (domain.UserModel, error)
}

type Service struct {
	userService    UserService
	sessionService service.SessionService
}

type Response struct {
	TokenPair
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(userService UserService, sessionService service.SessionService) Service {
	return Service{
		userService:    userService,
		sessionService: sessionService,
	}
}

// Login is a function that handles the login process
func (s Service) Login(request Request) (Response, error) {
	// Check if user exists
	exists := s.userService.ExistsByUsername(request.Username)
	if !exists {
		return Response{}, errors.New("user does not exist")
	}

	// Check if password is correct
	user, err := s.userService.VerifyUser(request.Username, request.Password)
	if err != nil {
		return Response{}, err
	}

	// Generate token
	tokenPair, err := GenerateTokenPair(user.ID)
	if err != nil {
		return Response{}, err
	}

	s.SetSession(user.ID.Hex(), tokenPair)

	return Response{
		tokenPair,
	}, nil
}

// SetSession is a function that sets the session with the given user id and token pair
func (s Service) SetSession(userId string, tokenPair TokenPair) error {
	// Start a session
	refreshTokenLifespan, err := strconv.Atoi(os.Getenv(RefreshExpirationKey))
	if err != nil {
		return err
	}

	// If there is a session for the user, delete it
	err = s.sessionService.DeleteSessionByUser(userId)
	if err != nil {
		return err
	}

	// Save refresh token with expire time
	createSessionCmd := commands.CreateSessionCommand{
		UserId:       userId,
		ExpireTime:   time.Now().Add(time.Hour * time.Duration(refreshTokenLifespan)).Unix(),
		RefreshToken: tokenPair.RefreshToken,
	}
	_, err = s.sessionService.CreateSession(createSessionCmd)
	if err != nil {
		return err
	}

	return nil
}

// GetCurrentUser is a function that returns the current user
func (s Service) GetCurrentUser(c *gin.Context) (domain.UserModel, error) {
	userId, err := ExtractUserIdFromContext(c)
	if err != nil {
		return domain.UserModel{}, err
	}

	user, err := s.userService.GetUser(userId)
	if err != nil {
		return domain.UserModel{}, err
	}

	return user, nil
}

// RefreshTokenPair is a function that refreshes the token pair
func (s Service) RefreshTokenPair(c *gin.Context) (TokenPair, error) {
	// Get the refresh token from the request body
	var refreshTokenRequest RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshTokenRequest); err != nil {
		return TokenPair{}, err
	}

	// Get current user id
	currentUser, err := s.GetCurrentUser(c)
	if err != nil {
		return TokenPair{}, err
	}

	// Get the session by user id
	sessionInfo, err := s.sessionService.GetUserSession(currentUser.ID.Hex())
	if err != nil {
		return TokenPair{}, err
	}

	// Check if the refresh token is valid
	if sessionInfo.RefreshToken != refreshTokenRequest.RefreshToken {
		return TokenPair{}, errors.New("refresh token is invalid")
	}

	// Now that we know the token is valid, we can extract the user id from it
	tokenPair, err := GenerateTokenPair(currentUser.ID)
	if err != nil {
		return TokenPair{}, err
	}

	s.SetSession(currentUser.ID.Hex(), tokenPair)

	return tokenPair, nil
}

// Logout is a function that logs out the user and invalidates the session
func (s Service) Logout(c *gin.Context) error {
	// Get current user
	currentUser, err := s.GetCurrentUser(c)
	if err != nil {
		return err
	}

	// Delete the session
	err = s.sessionService.DeleteSessionByUser(currentUser.ID.Hex())
	if err != nil {
		return err
	}

	return nil
}
