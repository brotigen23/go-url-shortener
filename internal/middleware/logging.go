package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Log(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//------------------------------------------------------------
			// REQUEST LOG
			//------------------------------------------------------------
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Errorln("failed to read request body:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			bodyCopy := bytes.NewBuffer(bodyBytes)
			logger.Debugln("Request Body:", bodyCopy.String())

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
			logger.Infoln(
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
}
