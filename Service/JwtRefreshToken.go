package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenRefreshToken(Login string) string {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": Login,
		"exp":      time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(key))
	if err != nil {
		return "failed"
	}
	return refreshTokenString
}
