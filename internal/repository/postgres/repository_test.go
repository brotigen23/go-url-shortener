package postgres

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brotigen23/go-url-shortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreate(t *testing.T) {
	shortURL := model.ShortURL{
		URL:      "google.com",
		ShortURL: "asd",
		Username: "user1",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	logger := zap.NewNop().Sugar()
	repo := New(db, logger)

	mock.ExpectBegin()
	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO short_url(url, short_url, username) VALUES($1, $2, $3)`)).
		WithArgs(
			shortURL.URL,
			shortURL.ShortURL,
			shortURL.Username,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Create(shortURL)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	rows := sqlmock.
		NewRows([]string{"id", "url", "short_url", "username", "is_deleted"}).
		AddRow(1, "123", "123", "user1", false)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	logger := zap.NewNop().Sugar()
	repo := New(db, logger)

	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT id, url, short_url, username, is_deleted FROM short_url`)).
		WithoutArgs().WillReturnRows(rows)

	_, err = repo.GetAll()
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetByUser(t *testing.T) {
	rows := sqlmock.
		NewRows([]string{"id", "url", "short_url", "is_deleted"}).
		AddRow(1, "123", "123", false)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	logger := zap.NewNop().Sugar()
	repo := New(db, logger)
	query := `
        SELECT id, url, short_url, is_deleted
        FROM short_url
        WHERE username = $1`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("user1").
		WillReturnRows(rows)

	_, err = repo.GetByUser("user1")
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetByURL(t *testing.T) {
	rows := sqlmock.
		NewRows([]string{"id", "short_url", "username", "is_deleted"}).
		AddRow(1, "1234", "user1", false)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	logger := zap.NewNop().Sugar()
	repo := New(db, logger)
	query := `
        SELECT id, short_url, username, is_deleted
        FROM short_url
        WHERE url = $1`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("google.com").
		WillReturnRows(rows)

	_, err = repo.GetByURL("google.com")
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetByAlias(t *testing.T) {
	rows := sqlmock.
		NewRows([]string{"id", "url", "username", "is_deleted"}).
		AddRow(1, "google.com", "user1", false)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	logger := zap.NewNop().Sugar()
	repo := New(db, logger)
	query := `
        SELECT id, url, username, is_deleted
        FROM short_url
        WHERE short_url = $1`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("1234").
		WillReturnRows(rows)

	_, err = repo.GetByAlias("1234")
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop().Sugar()
	repo := New(db, logger)

	query := `
        UPDATE short_url
        SET is_deleted = TRUE
        WHERE short_url IN ('1234')`

	mock.ExpectBegin()
	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithoutArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete("user1", []model.ShortURL{{ShortURL: "1234"}})
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetURLsCount(t *testing.T) {
	rows := sqlmock.
		NewRows([]string{"count"}).
		AddRow(1)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop().Sugar()
	repo := New(db, logger)

	query := `
	SELECT COUNT(*)
	FROM (SELECT DISTINCT url FROM short_url)
	AS temp`

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithoutArgs().
		WillReturnRows(rows)

	_ = repo.GetURLsCount()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUsersCount(t *testing.T) {
	rows := sqlmock.
		NewRows([]string{"count"}).
		AddRow(1)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop().Sugar()
	repo := New(db, logger)

	query := `
	SELECT COUNT(*)
	FROM (SELECT DISTINCT username FROM short_url)
	AS temp`

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithoutArgs().
		WillReturnRows(rows)

	_ = repo.GetUsersCount()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
