package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const target = "localhost:8080"

func TestAuth(t *testing.T) {
	logger := zap.NewNop().Sugar()
	middleware := New("secretKey", logger, "")
	username := utils.NewRandomString(16)

	expires := time.Hour * 1024
	jwtString, er := utils.BuildJWTString(username, "secretKey", expires)
	assert.NoError(t, er)
	cookie := &http.Cookie{
		Name:  "JWT",
		Value: jwtString,
	}
	emptyCookie := &http.Cookie{
		Name:  "JWT",
		Value: "",
	}

	type args struct {
		cookie *http.Cookie
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
			name: "Test No Cookie OK",
			args: args{
				cookie: nil,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test Has Cookie OK",
			args: args{
				cookie: cookie,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test Has Empty Cookie",
			args: args{
				cookie: emptyCookie,
			},
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, target, nil)
			if test.args.cookie != nil {
				request.AddCookie(test.args.cookie)
			}
			request.AddCookie(&http.Cookie{Name: "username", Value: "user"})

			w := httptest.NewRecorder()

			mid := middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			mid.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			if test.args.cookie == nil {
				coockies := result.Cookies()
				assert.NotEmpty(t, coockies)
			}
			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

func TestEncoding(t *testing.T) {
	logger := zap.NewNop().Sugar()
	middleware := New("secretKey", logger, "")
	type want struct {
		statusCode int
	}

	type args struct {
		body []byte
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				body: []byte("asd"),
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			_, err := gz.Write(test.args.body)
			assert.NoError(t, err)
			err = gz.Close()
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, target, &buf)
			request.Header.Set("Accept-Encoding", "gzip")
			request.Header.Set("Content-Encoding", "gzip")

			w := httptest.NewRecorder()

			mid := middleware.Encoding(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			mid.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

func TestLog(t *testing.T) {
	logger := zap.NewNop().Sugar()
	middleware := New("secretKey", logger, "")
	type want struct {
		statusCode int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "Test OK",
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, target, nil)

			w := httptest.NewRecorder()

			mid := middleware.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			mid.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

func TestSubnet(t *testing.T) {
	logger := zap.NewNop().Sugar()
	middleware := New("secretKey", logger, "127.0.0.1/16")

	type args struct {
		net string
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
				net: "127.0.0.14:23241",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test OK",
			args: args{
				net: "128.0.0.14:23241",
			},
			want: want{
				statusCode: http.StatusForbidden,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, target, nil)
			request.RemoteAddr = test.args.net
			w := httptest.NewRecorder()

			mid := middleware.CheckSubnet(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}))
			mid.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}
