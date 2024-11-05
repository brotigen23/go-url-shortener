package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type testGET struct {
	testName   string
	statusCode int
	model      *model.Alias
	location   string
}

func NewTest(testName string, statusCode int, model *model.Alias, location string) *testGET {
	t := &testGET{}
	t.testName = testName
	t.statusCode = statusCode
	t.model = model
	t.location = location
	return t
}

func TestIndexGetHandler(t *testing.T) {

	config := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "../../test/aliases.txt"}
	aliases, _ := utils.LoadLocalAliases(config.FileStoragePath)
	handler := NewIndexHandler(&config, aliases)

	tests := []*testGET{
		NewTest(
			"ya.ru test #1",
			http.StatusTemporaryRedirect,
			model.NewAlias("ya.ru", "VXcyQ01q"),
			"ya.ru"),
		NewTest(
			"google.com test #2",
			http.StatusTemporaryRedirect,
			model.NewAlias("yandex.ru", "K5IFupTM"),
			"yandex.ru"),
		NewTest(
			"not found test #3",
			http.StatusNotFound,
			model.NewAlias("", ""),
			""),
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+test.model.GetAlias(), nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.model.GetAlias())

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
