package jwt

import (
	"errors"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const tokenExp = time.Hour * 3

var ErrorInvalidJWT = errors.New("токен недействителен")

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuildJWT(userID, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString, secretKey string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга JWT токена: %w", err)
	}

	if !token.Valid {
		return "", ErrorInvalidJWT
	}

	logger.Log.Debug("Действительный JWT токен")

	return claims.UserID, nil
}
