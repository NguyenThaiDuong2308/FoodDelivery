package utils

import (
	"fmt"
	"shipper-service/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.LoadConfig().JWTSecret)

type Utils interface {
	ParseToken(token string) (jwt.MapClaims, error)
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
