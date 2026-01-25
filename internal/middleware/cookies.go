package middleware

import (
	"errors"
	"github.com/AleksandrTitov/shortener/internal/jwt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/userid"
	"net/http"
)

const idToken = "id_token"

func CookiesWrite(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			h.ServeHTTP(rw, r)
			return
		}
		var userID string
		var token string

		cooke, err := r.Cookie(idToken)
		if errors.Is(err, http.ErrNoCookie) {
			userID, _ = userid.New()
			logger.Log.Infof("Новый пользователь без токена, создаем новый User ID: %s", userID)

			token, err = jwt.BuildJWT(userID)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			logger.Log.Errorf("Ошибка получения cooke \"%s\": %v", idToken, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			token = cooke.Value
			userID, err = jwt.GetUserID(token)
			logger.Log.Infof("Пользователь с User ID: %s", userID)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		cookie := http.Cookie{
			Name:  idToken,
			Value: token,
		}

		http.SetCookie(rw, &cookie)
		h.ServeHTTP(rw, r)
	})
}
