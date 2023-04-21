package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const key = "12ZEFRGHJK4RT5YUJIKIOLIuytreds"

func GenJwtToken(Login string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": Login,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return string("err")
	}
	return tokenString
}

func GeneTokenFromRefreshToken(refreshTokenString string) string {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("invalid token:", err)
		return ""
	}
	username, ok := claims["username"].(string)
	if !ok {
		fmt.Println("missing or invalid username claim")
		return ""
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	signedAccessToken, err := accessToken.SignedString([]byte(key))
	if err != nil {
		fmt.Println("error signing access token:", err)
		return ""
	}
	return signedAccessToken
}
