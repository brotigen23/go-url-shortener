package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/mock"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const target = "http://localhost:8080"

// -------------------------------------------
// Init shared variables
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

	userService := service.New(cfg, logger, mockRepository)

	handler := New(cfg.BaseURL, userService)

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
		},
		{
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
			var err error
			switch test.args.contentType {
			case "text/plain":
				d = []byte(test.args.URL)
			case "application/json":
				d, err = json.Marshal(dto.ShortenRequest{URL: test.args.URL})
				assert.NoError(t, err)
			}

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

	userService := service.New(cfg, logger, mockRepository)

	handler := New(cfg.BaseURL, userService)

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
		mockRepository.EXPECT().Create(gomock.Any()).Return(nil),
		mockRepository.EXPECT().Create(gomock.Any()).Return(service.ErrShortURLAlreadyExists),
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

	userService := service.New(cfg, logger, mockRepository)

	handler := New(cfg.BaseURL, userService)

	type args struct {
		ShortURL string
	}
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test /{id} OK",
			args: args{
				ShortURL: "01234567",
			},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "ya.ru",
			},
		},
		{
			name: "Test /{id} Not Found",
			args: args{
				ShortURL: "123",
			},
			want: want{
				statusCode: http.StatusNotFound,
				location:   "",
			},
		},
		{
			name: "Test /{id} URL is Gone",
			args: args{
				ShortURL: "google.com",
			},
			want: want{
				statusCode: http.StatusGone,
				location:   "",
			},
		},
	}
	gomock.InOrder(
		// ------------------------------------------------------
		// Test /{id} OK
		// ------------------------------------------------------
		// Found url
		mockRepository.EXPECT().GetByAlias(gomock.Any()).Return(&model.ShortURL{URL: "ya.ru"}, nil),
		// Check is this deleted
		mockRepository.EXPECT().GetByAlias(gomock.Any()).Return(&model.ShortURL{URL: "ya.ru", IsDeleted: false}, nil),
		// ------------------------------------------------------
		// Test /{id} Not Found
		// ------------------------------------------------------
		mockRepository.EXPECT().GetByAlias(gomock.Any()).Return(nil, service.ErrShortURLNotFound),

		// ------------------------------------------------------
		// Test /{id} URL is Gone
		// ------------------------------------------------------
		mockRepository.EXPECT().GetByAlias(gomock.Any()).Return(&model.ShortURL{}, nil),
		mockRepository.EXPECT().GetByAlias(gomock.Any()).Return(&model.ShortURL{IsDeleted: true}, nil),
	)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			d := []byte(test.args.ShortURL)
			request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(d))
			request.AddCookie(&http.Cookie{Name: "username", Value: "user"})
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.args.ShortURL)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			handler.RedirectByShortURL(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			if test.want.statusCode == http.StatusTemporaryRedirect {
				location := result.Header.Get("location")
				assert.Equal(t, test.want.location, location)
			}
		})
	}
}
