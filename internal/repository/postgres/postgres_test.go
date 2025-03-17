package postgres

import (
	"database/sql"
	"os"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	/////////////////////////////////////////////////////
	// CHECK CONNECTION
	/////////////////////////////////////////////////////
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.SkipNow()
	}

	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	assert.NoError(t, err)

	logger := zap.NewNop()

	assert.NoError(t, err)
	repo := New(db, logger.Sugar())

	/////////////////////////////////////////////////////
	// CREATE
	/////////////////////////////////////////////////////
	shortURLsToSave := []model.ShortURL{
		{URL: "ya.ru", ShortURL: "sadfasdf3", Username: "user1"},
		{URL: "google.com", ShortURL: "sadf23sdav", Username: "user1"},
		{URL: "metanit.com", ShortURL: "sdaf4423d", Username: "user2"},
	}
	// OK
	err = repo.Create(shortURLsToSave[0])
	assert.NoError(t, err)
	err = repo.Create(shortURLsToSave[1])
	assert.NoError(t, err)
	err = repo.Create(shortURLsToSave[2])
	assert.NoError(t, err)
	// ERR
	err = repo.Create(shortURLsToSave[0])
	assert.Equal(t, repository.ErrShortURLAlreadyExists, err)

	/////////////////////////////////////////////////////
	// GET
	/////////////////////////////////////////////////////
	_, err = repo.GetAll()
	assert.NoError(t, err)

	_, err = repo.GetByUser("user1")
	assert.NoError(t, err)

	_, err = repo.GetByURL("google.com")
	assert.NoError(t, err)

	_, err = repo.GetByAlias("sadfasdf3")
	assert.NoError(t, err)

	err = repo.Delete("user1", []model.ShortURL{
		shortURLsToSave[0],
	})
	assert.NoError(t, err)

	userCount := repo.GetUsersCount()
	assert.Equal(t, 2, userCount)

	urlsCount := repo.GetURLsCount()
	assert.Equal(t, 3, urlsCount)
}
