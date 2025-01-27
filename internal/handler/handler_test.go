package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/mock"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const target = "http://localhost:8080"

func TestCreateShortURL(t *testing.T) {
	config, err := config.NewConfig()
	require.NoError(t, err)
	logger := zap.NewNop().Sugar()
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepository := mock.NewMockRepository(controller)

	userService, err := service.New(config, logger, mockRepository)
	require.NoError(t, err)

	handler, err := New(config.BaseURL, userService)
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
	config, err := config.NewConfig()
	require.NoError(t, err)
	logger := zap.NewNop().Sugar()
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepository := mock.NewMockRepository(controller)

	userService, err := service.New(config, logger, mockRepository)
	require.NoError(t, err)

	handler, err := New(config.BaseURL, userService)
	require.NoError(t, err)

	type args struct {
		data []dto.BatchRequest
	}
	type want struct {
		response   []dto.BatchRequest
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
				data: []dto.BatchRequest{
					{
						ID:  "0",
						URL: "example.com",
					},
					{
						ID:  "1",
						URL: "example.com",
					},
				},
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		}, {
			name: "Test Conflict text",
			args: args{
				data:        "google.com",
				contentType: "text/plain",
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "Test Incorrect data text",
			args: args{
				data:        "",
				contentType: "text/plain",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Test OK json",
			args: args{
				data:        "ya.ru",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		}, {
			name: "Test Conflict json",
			args: args{
				data:        "google.com",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "Test Incorrect data json",
			args: args{
				data:        "",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
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
				d = []byte(test.args.data)
			case "application/json":
				d, err = json.Marshal(dto.ShortenRequest{URL: test.args.data})
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
