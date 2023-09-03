package commands

type CreateSessionCommand struct {
	UserId       string `json:"userId"`
	ExpireTime   int64  `json:"expireTime"`
	RefreshToken string `json:"refreshToken"`
}

type DeleteSessionCommand struct {
	UserId string `json:"userId"`
}
