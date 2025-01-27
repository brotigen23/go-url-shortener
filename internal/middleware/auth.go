package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/utils"
	"go.uber.org/zap"
)

func Auth(key string, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("JWT")

			if err != nil {
				if err == http.ErrNoCookie {
					username := utils.NewRandomString(16)

					expires := time.Hour * 1024
					jwtString, err := utils.BuildJWTString(username, key, expires)
					if err != nil {
						log.Printf("error: %v", err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					cookie = &http.Cookie{
						Name:  "JWT",
						Value: jwtString,
					}
					http.SetCookie(w, cookie)
					r.AddCookie(&http.Cookie{Name: "username", Value: username})
					next.ServeHTTP(w, r)
					return
				} else {
					logger.Errorf("auth error: %n\n", err.Error())
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return

				}
			}

			user, err := utils.GetUsernameFromJWT(cookie.Value, key)
			if err != nil {
				logger.Errorf("auth error: %n\n", err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			r.AddCookie(&http.Cookie{Name: "username", Value: user})
			next.ServeHTTP(w, r)
		})
	}
}
