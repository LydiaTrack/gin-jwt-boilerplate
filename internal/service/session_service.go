package service

import (
	"errors"
	"gin-jwt-boilerplate/internal/domain"
	"gin-jwt-boilerplate/internal/domain/commands"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// SessionService is an interface that contains the methods for the session service
type SessionService struct {
	sessionRepository SessionRepository
	UserService
}

// SessionRepository is an interface that contains the methods for the session repository
type SessionRepository interface {
	// SaveSession is a function that creates a session
	SaveSession(model domain.SessionInfoModel) (domain.SessionInfoModel, error)
	// GetUserSession is a function that gets a user session
	GetUserSession(id bson.ObjectId) (domain.SessionInfoModel, error)
	// DeleteSessionByUserId is a function that deletes a session
	DeleteSessionByUserId(id bson.ObjectId) error
	// DeleteSessionById is a function that deletes a session by id
	DeleteSessionById(sessionId bson.ObjectId) error
}

func NewSessionService(sessionRepository SessionRepository, userService UserService) SessionService {
	return SessionService{
		sessionRepository: sessionRepository,
		UserService:       userService,
	}
}

// CreateSession is a function that creates a session
func (s SessionService) CreateSession(cmd commands.CreateSessionCommand) (domain.SessionInfoModel, error) {
	// Check if user exists
	exists, err := s.UserService.ExistsUser(cmd.UserId)
	if err != nil {
		return domain.SessionInfoModel{}, err
	}
	if !exists {
		return domain.SessionInfoModel{}, errors.New("user does not exist")
	}
	sessionInfo := domain.SessionInfoModel{
		ID:           bson.NewObjectId(),
		UserId:       bson.ObjectIdHex(cmd.UserId),
		ExpireTime:   cmd.ExpireTime,
		RefreshToken: cmd.RefreshToken,
	}
	// TODO: Permission check
	return s.sessionRepository.SaveSession(sessionInfo)
}

// GetUserSession is a function that gets a user session
func (s SessionService) GetUserSession(id string) (domain.SessionInfoModel, error) {
	// Check if user exists
	exists, err := s.UserService.ExistsUser(id)
	if err != nil {
		return domain.SessionInfoModel{}, err
	}
	if !exists {
		return domain.SessionInfoModel{}, errors.New("user does not exist")
	}

	return s.sessionRepository.GetUserSession(bson.ObjectIdHex(id))
}

// DeleteSession is a function that deletes a session
func (s SessionService) DeleteSessionByUser(userId string) error {
	return s.sessionRepository.DeleteSessionByUserId(bson.ObjectIdHex(userId))
}

// DeleteSessionById is a function that deletes a session by id
func (s SessionService) DeleteSessionById(sessionId string) error {
	return s.sessionRepository.DeleteSessionById(bson.ObjectIdHex(sessionId))
}

// IsUserHasActiveSession is a function that checks if a user has an active session
func (s SessionService) IsUserHasActiveSession(userId string) bool {
	session, err := s.GetUserSession(userId)
	if err != nil {
		return false
	}

	// Check if session still valid by comparing the expire time with the current time
	currentTime := time.Now().Unix()
	if session.ExpireTime < currentTime {
		return false
	}

	return true
}
