package handler

import (
	"net/http"
	"testing"
)

func TestSaveURL(t *testing.T) {
	type args struct {
		data        string
		contentType string
	}
	type want struct {
		statusCode int
	}
	_ = []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				data:        "asd",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		}, {
			name: "Test Conflict",
			args: args{
				data:        "",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		}, {
			name: "Test Incorrect data",
			args: args{
				data:        "123",
				contentType: "application/json",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

}
