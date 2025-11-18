package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func ValidateJWT(tokenString string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		//TODO: check the algorithm
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("unauthenticated")
	}

	claims := token.Claims.(*jwt.StandardClaims)
	id := claims.Issuer
	return id, nil
}
