package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))
var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

func GenerateAccessToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}
