package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/mock"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const target = "http://localhost:8080"

// -------------------------------------------
// Init shared var
// -------------------------------------------
var cfg = &config.Config{
	ServerAddress:   "localhost:8080",
	BaseURL:         "http://localhost:8080",
	FileStoragePath: "",
	DatabaseDSN:     "",

	JWTSecretKey: "secret",
}
var logger = zap.NewNop().Sugar()

func TestCreateShortURL(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepository := mock.NewMockRepository(controller)

	userService, err := service.New(cfg, logger, mockRepository)
	require.NoError(t, err)

	handler, err := New(cfg.BaseURL, userService)
	require.NoError(t, err)

	type args struct {
		URL         string
		contentType string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK text",
			args: args{
				URL:         "ya.ru",
				contentType: "text/plain",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		}, {
			name: "Test Conflict text",
			args: args{
				URL:         "google.com",
				contentType: "text/plain",
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "Test Incorrect data text",
			args: args{
				URL:         "",
				contentType: "text/plain",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Test OK json",
			args: args{
				URL:         "ya.ru",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		}, {
			name: "Test Conflict json",
			args: args{
				URL:         "google.com",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "Test Incorrect data json",
			args: args{
				URL:         "",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	mockRepository.EXPECT().GetByURL(gomock.Any()).Return(&model.ShortURL{}, nil).MaxTimes(2)
	gomock.InOrder(
		mockRepository.EXPECT().Create(gomock.Any()).Return(nil),
		mockRepository.EXPECT().Create(gomock.Any()).Return(repository.ErrShortURLAlreadyExists),
		mockRepository.EXPECT().Create(gomock.Any()).Return(nil),
		mockRepository.EXPECT().Create(gomock.Any()).Return(repository.ErrShortURLAlreadyExists),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var d []byte
			switch test.args.contentType {
			case "text/plain":
				d = []byte(test.args.URL)
			case "application/json":
				d, err = json.Marshal(dto.ShortenRequest{URL: test.args.URL})
			}
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(d))
			request.Header.Set("content-type", test.args.contentType)
			request.AddCookie(&http.Cookie{Name: "username", Value: "user"})
			w := httptest.NewRecorder()

			handler.CreateShortURL(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

func TestCreateShortURLs(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepository := mock.NewMockRepository(controller)

	userService, err := service.New(cfg, logger, mockRepository)
	require.NoError(t, err)

	handler, err := New(cfg.BaseURL, userService)
	require.NoError(t, err)

	type args struct {
		URLs []dto.BatchRequest
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				URLs: []dto.BatchRequest{
					{
						ID:  "0",
						URL: "ya.ru",
					},
					{
						ID:  "1",
						URL: "google.com",
					},
				},
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		}, {
			name: "Test Conflict",
			args: args{
				URLs: []dto.BatchRequest{
					{
						ID:  "0",
						URL: "ya.ru",
					},
					{
						ID:  "1",
						URL: "google.com",
					},
				},
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "Test Incorrect Data",
			args: args{
				URLs: []dto.BatchRequest{
					{
						ID:  "0",
						URL: "",
					},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	mockRepository.EXPECT().GetByURL(gomock.Any()).Return(&model.ShortURL{}, nil).MaxTimes(2)
	gomock.InOrder(
		mockRepository.EXPECT().Create(gomock.Any()).Return(nil),
		mockRepository.EXPECT().Create(gomock.Any()).Return(repository.ErrShortURLAlreadyExists),
		mockRepository.EXPECT().Create(gomock.Any()).Return(nil),
		mockRepository.EXPECT().Create(gomock.Any()).Return(repository.ErrShortURLAlreadyExists),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, err := json.Marshal(test.args.URLs)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(d))
			request.Header.Set("content-type", "application/json")
			request.AddCookie(&http.Cookie{Name: "username", Value: "user"})
			w := httptest.NewRecorder()

			handler.CreateShortURLs(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

func TestRedirectByShortURL(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepository := mock.NewMockRepository(controller)

	userService, err := service.New(cfg, logger, mockRepository)
	require.NoError(t, err)

	handler, err := New(cfg.BaseURL, userService)
	require.NoError(t, err)

	regexCorrectHeaderLocation, err := regexp.Compile("http://" + cfg.BaseURL + "/" + "\\w{" + "8" + "}")
	require.NoError(t, err)

	type args struct {
		ShortURL string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test /{id} OK",
			args: args{
				ShortURL: "ya.ru",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "Test /{id} Incorrect URL",
			args: args{
				ShortURL: "",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Test /api/shorten/batch OK",
			args: args{
				ShortURL: "ya.ru",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		}, {
			name: "Test /api/shorten/batch Incorrect Data",
			args: args{
				ShortURL: "google.com",
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
	}
	mockRepository.EXPECT().GetByURL(gomock.Any()).Return(&model.ShortURL{}, nil).MaxTimes(2)
	gomock.InOrder(
		mockRepository.EXPECT().Create(gomock.Any()).Return(nil),
		mockRepository.EXPECT().Create(gomock.Any()).Return(repository.ErrShortURLAlreadyExists),
		mockRepository.EXPECT().Create(gomock.Any()).Return(nil),
		mockRepository.EXPECT().Create(gomock.Any()).Return(repository.ErrShortURLAlreadyExists),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			d := []byte(test.args.ShortURL)
			request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(d))

			request.AddCookie(&http.Cookie{Name: "username", Value: "user"})
			w := httptest.NewRecorder()

			handler.RedirectByShortURL(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			if test.want.statusCode == http.StatusTemporaryRedirect {
				location := result.Header.Get("location")
				assert.Regexp(t, regexCorrectHeaderLocation, location)
			}
		})
	}
}
