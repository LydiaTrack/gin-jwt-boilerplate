package auth

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	UserIdKey            = "sub"
	AuthorizedKey        = "authorized"
	ExpKey               = "exp"
	JwtExpirationKey     = "JWT_EXPIRES_IN_HOUR"
	RefreshExpirationKey = "JWT_REFRESH_EXPIRES_IN_HOUR"
	JwtSecretKey         = "JWT_SECRET"
	AuthorizationHeader  = "Authorization"
)

type TokenPair struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// GenerateTokenPair generates a jwt and refresh token
func GenerateTokenPair(userId bson.ObjectId) (TokenPair, error) {

	tokenLifespan, err := strconv.Atoi(os.Getenv(JwtExpirationKey))

	if err != nil {
		return TokenPair{}, err
	}

	claims := jwt.MapClaims{}
	claims[AuthorizedKey] = true
	claims[UserIdKey] = userId.Hex()
	claims[ExpKey] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(os.Getenv(JwtSecretKey)))
	if err != nil {
		return TokenPair{}, err
	}

	// Refresh token is a random string
	refreshTokenStr := bson.NewObjectId().Hex()

	return TokenPair{Token: tokenStr, RefreshToken: refreshTokenStr}, nil
}

// IsTokenValid validates the token
func IsTokenValid(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(JwtSecretKey)), nil
	})

	if err != nil {
		return err
	}
	return nil
}

// ExtractTokenFromContext extracts the user id from the bearer token or refresh token
func ExtractTokenFromContext(c *gin.Context) (string, error) {
	token := extractBearerTokenFromContext(c)

	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

// extractBearerTokenFromContext extracts the token from the request
func extractBearerTokenFromContext(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get(AuthorizationHeader)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractUserIdFromContext extracts the token id (userId) from the request
func ExtractUserIdFromContext(c *gin.Context) (string, error) {

	tokenString, err := ExtractTokenFromContext(c)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(JwtSecretKey)), nil
	})

	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid := claims[UserIdKey].(string)
		return uid, nil
	}
	return "", nil
}
