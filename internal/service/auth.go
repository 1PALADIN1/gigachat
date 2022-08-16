package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const defaultUser = "username" //TODO: from DB

func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": defaultUser,
		"iat": time.Now().UTC(),
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
