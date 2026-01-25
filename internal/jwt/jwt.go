package jwt

import (
	"errors"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

var ErrorInvalidJWT = errors.New("токен недействителен")

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuildJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		logger.Log.Error("Недействительный JWT токен")
		return "", ErrorInvalidJWT
	}

	logger.Log.Debug("Действительный JWT токен")

	return claims.UserID, nil
}
