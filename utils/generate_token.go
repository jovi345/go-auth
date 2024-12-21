package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID string, email string, firstName string, lastName string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET_ACCESS")
	expirationSeconds := 15

	claims := jwt.MapClaims{
		"user_id":    userID,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"role":       role,
		"exp":        time.Now().Add(time.Duration(expirationSeconds) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID string, email string, firstName string, lastName string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET_REFRESH")
	expirationDays := 5

	claims := jwt.MapClaims{
		"user_id":    userID,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"role":       role,
		"exp":        time.Now().Add(time.Duration(expirationDays*24) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
