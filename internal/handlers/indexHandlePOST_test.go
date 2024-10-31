package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexHandePOST(t *testing.T) {

	config := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	handler := NewIndexHandler(&config)

	responseRegexp, _ := regexp.Compile("http://" + config.ServerAddress + "/" + "\\w{" + "8" + "}")

	type want struct {
		statusCode  int
		contentType string
	}

	tests := []struct {
		testName string
		url      string
		want     want
	}{
		{
			testName: "test #1",
			url:      "https://ya.ru",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			testName: "test #2",
			url:      "https://yandex.ru",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()
			handler.HandlePOST(w, request)
			result := w.Result()

			// получаем и проверяем тело запроса
			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)

			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("content-type"))
			assert.Regexp(t, responseRegexp, string(resBody))
		})
	}

}
