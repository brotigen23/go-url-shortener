package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexHandlePOST(t *testing.T) {

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
			url:      "https://123.ru",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			testName: "test #2",
			url:      "https://1232.ru",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(test.url)))
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

func TestIndexHandlePOSTAPI(t *testing.T) {

	config := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	handler := NewIndexHandler(&config)

	//responseRegexp, _ := regexp.Compile("http://" + config.ServerAddress + "/" + "\\w{" + "8" + "}")

	type want struct {
		statusCode  int
		contentType string
		resp        resp
	}

	tests := []struct {
		testName string
		url      req
		want     want
	}{
		{
			testName: "test #1",
			url:      req{"https://ya.ru"},
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				resp:        resp{"asd"},
			},
		},
		{
			testName: "test #2",
			url:      req{"https://yandex.ru"},
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				resp:        resp{"asd"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			body, _ := json.Marshal(test.url)
			request, _ := http.NewRequest("POST", "/api/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()
			handler.HandlePOSTAPI(w, request)
			result := w.Result()

			// получаем и проверяем тело запроса
			defer result.Body.Close()
			_, err := io.ReadAll(result.Body)

			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("content-type"))
		})
	}

}
