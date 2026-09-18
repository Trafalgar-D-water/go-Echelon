package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessTokenSecret  []byte
	refreshTokenSecret []byte
)

type TokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func ConfigureJWT(accessSecret, refreshSecret string) error {
	if len(accessSecret) < 32 {
		return errors.New("ACCESS_TOKEN_SECRET must be at least 32 characters")
	}

	if len(refreshSecret) < 32 {
		return errors.New("REFRESH_TOKEN_SECRET must be at least 32 characters")
	}

	if accessSecret == refreshSecret {
		return errors.New("access and refresh token secrets must be different")
	}

	accessTokenSecret = []byte(accessSecret)
	refreshTokenSecret = []byte(refreshSecret)

	return nil
}

func GenerateAccessToken(userID string) (string, error) {
	return generateToken(userID, 15*time.Minute, accessTokenSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	return generateToken(userID, 7*24*time.Hour, refreshTokenSecret)
}

func ParseAccessToken(tokenString string) (*TokenClaims, error) {
	return parseToken(tokenString, accessTokenSecret)
}

func ParseRefreshToken(tokenString string) (*TokenClaims, error) {
	return parseToken(tokenString, refreshTokenSecret)
}

func generateToken(userID string, duration time.Duration, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("JWT secrets have not been configured")
	}

	now := time.Now().UTC()

	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func parseToken(tokenString string, secret []byte) (*TokenClaims, error) {
	if len(secret) == 0 {
		return nil, errors.New("JWT secrets have not been configured")
	}

	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.UserID == "" {
		return nil, errors.New("token does not contain a user ID")
	}

	return claims, nil
}
