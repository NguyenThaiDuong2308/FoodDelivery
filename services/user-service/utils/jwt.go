package utils

import (
	"errors"
	"fmt"
	"time"
	"user-service/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.LoadConfig().JWTSecret)

type Utils interface {
	GenerateToken(userID uint, email string, role string) (string, error)
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

func GenerateToken(userID uint, role string, email string, expiryTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(expiryTime).Unix(),
		"role":    role,
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
