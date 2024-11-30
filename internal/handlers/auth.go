package handlers

import (
	"fmt"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/services"
)

func WithAuth(next http.HandlerFunc, config *config.Config, service *services.ServiceAuth) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Считываем значение
			userID, err := r.Cookie("userID")
			switch err {
			case http.ErrNoCookie:
				// Генерируем нового пользователя
				id, err := service.GenerateID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				// Сохраняем созданного пользователя
				service.SaveUser(id)
				http.SetCookie(w, &http.Cookie{
					Name:  "userID",
					Value: id,
				})
				r.AddCookie(&http.Cookie{
					Name:  "userID",
					Value: id,
				})
				fmt.Println(r.Cookie("userID"))
				// Подписываем нового пользователя
				err = service.SignUser(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case nil:
				if userID.Value == "" {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				// Если значение есть проверяем подпись
				if !service.CheckSing(userID.Value) {
					// Если подпись недействительная то генерируем новую
					err = service.SignUser(userID.Value)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
				}
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			next.ServeHTTP(w, r)
		})
}