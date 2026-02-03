package middleware

import (
	"errors"
	"github.com/AleksandrTitov/shortener/internal/jwt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/userid"
	"net/http"
)

const (
	idToken      = "id_token"
	UserIDHeader = "X-User-ID"
)

func CookiesJWT(secretKey string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var userID string
			var token string

			cooke, err := r.Cookie(idToken)
			if errors.Is(err, http.ErrNoCookie) {
				userID, err = userid.New()
				if err != nil {
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				logger.Log.Infof("Новый пользователь без токена, создаем новый User ID: %s", userID)

				token, err = jwt.BuildJWT(userID, secretKey)
				if err != nil {
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else if err != nil {
				logger.Log.Debugf("Ошибка получения cooke %q: %v", idToken, err)
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			} else {
				token = cooke.Value
				userID, err = jwt.GetUserID(token, secretKey)
				if err != nil {
					logger.Log.Debugf("Ошибка получения User ID из JWT: %v", err)
					http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
				logger.Log.Infof("Пользователь с User ID: %s", userID)
			}

			cookie := http.Cookie{
				Name:  idToken,
				Value: token,
			}

			r.Header.Set(UserIDHeader, userID)

			http.SetCookie(rw, &cookie)
			h.ServeHTTP(rw, r)
		})
	}
}
