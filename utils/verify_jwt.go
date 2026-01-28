package utils

import (
	"errors"
	"money/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string) (string, error) {
	configuration := config.GetConfig()

	jwtSecret := []byte(configuration.JwtSecret)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (string, error) {
	configuration := config.GetConfig()

	jwtSecret := []byte(configuration.JwtSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)

		if !ok {
			return "", errors.New("invalid token")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
