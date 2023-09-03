package commands

import (
	"gin-jwt-boilerplate/internal/domain"
	"gopkg.in/mgo.v2/bson"
)

type CreateUserCommand struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	domain.PersonInfo `json:"person_info"`
}

type UpdateUserCommand struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	domain.PersonInfo `json:"person_info"`
}

type DeleteUserCommand struct {
	ID bson.ObjectId `json:"_id"`
}
