package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestPOST(t *testing.T) {
	type want struct {
		code        int
		URL         string
		request     string
		aliasPrefix string
		aliasLenght int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test #1. ya.ru",
			want: want{
				code:        201,
				URL:         `https://ya.ru`,
				request:     "http://localhost:8080/",
				aliasPrefix: "http://localhost:8080",
				aliasLenght: storage.ALIASLENGHT,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.want.request, strings.NewReader(test.want.URL))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			IndexHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
func TestGET(t *testing.T) {
	type want struct {
		code     int
		location string
		request  string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test #1. ya.ru",
			want: want{
				code:     307,
				location: `https://ya.ru`,
				request:  "http://localhost:8080/asdfgh",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.want.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			IndexHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			assert.Equal(t, test.want.location, string(res.Header.Get("location")))
		})
	}
}
