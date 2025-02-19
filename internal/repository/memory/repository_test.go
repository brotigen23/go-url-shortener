package memory

import (
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortURL(t *testing.T) {
	repo := New([]model.ShortURL{
		{ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
		{ID: 2, URL: "asd", ShortURL: "asd", Username: "asd", IsDeleted: false},
	})

	type args struct {
		ShortURL model.ShortURL
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				ShortURL: model.ShortURL{URL: "url1", ShortURL: "shorturl", Username: "user", IsDeleted: false},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Test Already Exists",
			args: args{
				ShortURL: model.ShortURL{URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
			},
			want: want{
				err: repository.ErrShortURLAlreadyExists,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Create(test.args.ShortURL)
			assert.ErrorIs(t, err, test.want.err)
		})
	}
}

func TestGetByUser(t *testing.T) {
	repo := New([]model.ShortURL{
		{ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
		{ID: 2, URL: "asd", ShortURL: "asd", Username: "user", IsDeleted: false},
		{ID: 3, URL: "asd", ShortURL: "asd", Username: "asd", IsDeleted: false},
	})

	type args struct {
		username string
	}
	type want struct {
		shortURLs []model.ShortURL
		err       error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				username: "user",
			},
			want: want{
				shortURLs: []model.ShortURL{
					{ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
					{ID: 2, URL: "asd", ShortURL: "asd", Username: "user", IsDeleted: false},
				},
				err: nil,
			},
		},
		{
			name: "Test Not Found",
			args: args{
				username: "user1",
			},
			want: want{
				shortURLs: nil,
				err:       repository.ErrNoFound,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls, err := repo.GetByUser(test.args.username)
			assert.ErrorIs(t, err, test.want.err)
			assert.Equal(t, test.want.shortURLs, urls)
		})
	}

}

func TestGetByURL(t *testing.T) {
	repo := New([]model.ShortURL{
		{ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
		{ID: 2, URL: "ya.ru", ShortURL: "qwerty", Username: "user", IsDeleted: false},
		{ID: 3, URL: "example.org", ShortURL: "asd", Username: "asd", IsDeleted: false},
	})

	type args struct {
		url string
	}
	type want struct {
		shortURLs *model.ShortURL
		err       error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				url: "google.com",
			},
			want: want{
				shortURLs: &model.ShortURL{
					ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false,
				},
				err: nil,
			},
		},
		{
			name: "Test Not Found",
			args: args{
				url: "someurl",
			},
			want: want{
				shortURLs: nil,
				err:       repository.ErrNoFound,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls, err := repo.GetByURL(test.args.url)
			assert.ErrorIs(t, err, test.want.err)
			assert.Equal(t, test.want.shortURLs, urls)
		})
	}
}

func TestGetByAlias(t *testing.T) {
	repo := New([]model.ShortURL{
		{ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
		{ID: 2, URL: "ya.ru", ShortURL: "qwerty", Username: "user", IsDeleted: false},
		{ID: 3, URL: "example.org", ShortURL: "asd", Username: "asd", IsDeleted: false},
	})

	type args struct {
		alias string
	}
	type want struct {
		shortURLs *model.ShortURL
		err       error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				alias: "12345678",
			},
			want: want{
				shortURLs: &model.ShortURL{
					ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false,
				},
				err: nil,
			},
		},
		{
			name: "Test Not Found",
			args: args{
				alias: "somealias",
			},
			want: want{
				shortURLs: nil,
				err:       repository.ErrNoFound,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls, err := repo.GetByAlias(test.args.alias)
			assert.ErrorIs(t, err, test.want.err)
			assert.Equal(t, test.want.shortURLs, urls)
		})
	}
}

func TestDelete(t *testing.T) {
	repo := New([]model.ShortURL{
		{ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
		{ID: 2, URL: "ya.ru", ShortURL: "qwerty", Username: "user", IsDeleted: false},
		{ID: 3, URL: "example.org", ShortURL: "asd", Username: "asd", IsDeleted: false},
	})

	type args struct {
		shortURL []model.ShortURL
		username string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test OK",
			args: args{
				shortURL: []model.ShortURL{
					{ID: 1, URL: "google.com", ShortURL: "12345678", Username: "user", IsDeleted: false},
				},
				username: "user",
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Delete(test.args.username, test.args.shortURL)
			assert.ErrorIs(t, err, test.want.err)
		})
	}
}
