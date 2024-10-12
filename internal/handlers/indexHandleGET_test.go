package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var stor = storage.NewStorage()

type testGET struct {
	testName   string
	statusCode int
	url        string
	alias      string
	location   string
}

func NewTest(testName string, statusCode int, url string, location string) *testGET {
	t := new(testGET)
	t.testName = testName
	t.statusCode = statusCode
	t.url = url
	t.location = location
	t.alias, _ = stor.FindByURL([]byte(url))
	return t
}

func TestIndexGetHandler(t *testing.T) {
	stor.Put([]byte("https://ya.ru"))
	stor.Put([]byte("https://google.com"))
	stor.Put([]byte("https://rutube.ru"))
	stor.Put([]byte("http://metanit.com"))

	tests := []*testGET{
		NewTest(
			"ya.ru test #1",
			http.StatusTemporaryRedirect,
			"https://ya.ru",
			"https://ya.ru"),
		NewTest(
			"google.com test #2",
			http.StatusTemporaryRedirect,
			"https://google.com",
			"https://google.com"),
		NewTest(
			"not found test #3",
			http.StatusNotFound,
			"https://yandex.ru",
			""),
	}

	config := config.NewConfig()
	handler := NewIndexHandler(config, stor)

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+test.alias, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.alias)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			handler.HandleGET(w, request)
			result := w.Result()
			defer result.Body.Close()

			// Status code
			assert.Equal(t, test.statusCode, result.StatusCode)
			assert.Equal(t, test.location, result.Header.Get("location"))

		})
	}
}
