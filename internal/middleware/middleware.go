package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/utils"
	"go.uber.org/zap"
)

// Структура Middleware
type Middleware struct {
	logger    *zap.SugaredLogger
	secretKey string
}

// Создание экземпляра Middleware
func New(secretKey string, logger *zap.SugaredLogger) *Middleware {
	return &Middleware{
		logger:    logger,
		secretKey: secretKey,
	}
}

// Middleware для аутентификации пользователя
func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("JWT")
		if err != nil {
			if err == http.ErrNoCookie {
				if r.URL.Path == "/api/user/urls" {
					m.logger.Errorln("No cookie")
					w.WriteHeader(http.StatusNoContent)
					return
				}
				username := utils.NewRandomString(16)
				m.logger.Debugln("new user", username)

				expires := time.Hour * 1024
				jwtString, er := utils.BuildJWTString(username, m.secretKey, expires)
				if er != nil {
					m.logger.Errorln(er)
					http.Error(w, er.Error(), http.StatusInternalServerError)
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
				m.logger.Errorln(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}
		if cookie.Value == "" {
			m.logger.Errorln("Coockie is empty")
			http.Error(w, "username is empty", http.StatusUnauthorized)
			return
		}
		username, err := utils.GetUsernameFromJWT(cookie.Value, m.secretKey)
		if err != nil {
			m.logger.Errorln(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r.AddCookie(&http.Cookie{Name: "username", Value: username})
		m.logger.Debug("request from user: ", username)
		next.ServeHTTP(w, r)
	})
}

// Middleware для сжатия ответа
func (m *Middleware) Encoding(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

// Производит логгирование входящего запроса
func (m *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//------------------------------------------------------------
		// REQUEST LOG
		//------------------------------------------------------------
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			m.logger.Errorln("failed to read request body:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		bodyCopy := bytes.NewBuffer(bodyBytes)
		m.logger.Debugln("Request Body:", bodyCopy.String())

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		//------------------------------------------------------------
		// RESPONSE LOG
		//------------------------------------------------------------

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)
		m.logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
			"encoding", r.Header.Get("Content-Encoding"),
			"content-type", r.Header.Get("content-type"),
		)
	})
}
