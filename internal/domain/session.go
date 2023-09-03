package domain

import "gopkg.in/mgo.v2/bson"

// SessionInfoModel is a struct that contains the session information and maps to the userId
type SessionInfoModel struct {
	ID           bson.ObjectId `bson:"_id"`
	UserId       bson.ObjectId `bson:"userId"`
	ExpireTime   int64         `bson:"expireTime"`
	RefreshToken string        `bson:"refreshToken"`
}
