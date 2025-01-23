package middleware

import (
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"go.uber.org/zap"
)

func Auth(config *config.Config, service *service.ServiceAuth, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Считываем значение
			userID, err := r.Cookie("userID")
			switch err {
			case http.ErrNoCookie:
				// Генерируем нового пользователя
				userName, err := service.GenerateID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					logger.Errorln("error creating new id")
					return
				}
				// Сохраняем созданного пользователя
				err = service.SaveUser(userName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					logger.Errorln("error to save new user")
					return
				}
				logger.Infoln("new user:", userName, "saved")
				http.SetCookie(w, &http.Cookie{
					Name:  "userID",
					Value: userName,
					Path:  "/api/user/urls",
				})
				r.AddCookie(&http.Cookie{
					Name:  "userID",
					Value: userName,
					Path:  "/api/user/urls",
				})
				// Подписываем нового пользователя
				err = service.SignUser(userName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if r.URL.Path == "/api/user/urls" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			case nil:
				if userID.Value == "" {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					logger.Errorln("no username in cookie")
					return
				}
				// Если значение есть проверяем подпись
				if !service.CheckSing(userID.Value) {
					logger.Warnln("sign is invalid")
					// Если подпись недействительная то генерируем новую
					err = service.SignUser(userID.Value)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
				}
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
