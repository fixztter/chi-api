package auth

import (
	"strconv"
	"time"

	"github.com/fixztter/chi-api/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userID int) (string, error) {
	expirationTime, _ := strconv.ParseInt(config.Envs.JWTExpirationSeconds, 10, 64)
	expiration := time.Second * time.Duration(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    strconv.Itoa(userID),
		"expired_at": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
