package repository

import (
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const pgDriver = "postgres"
const stringConnection = "host=localhost port=5432 user=myuser password=1234 dbname=mydb sslmode=disable"

func TestConnection(t *testing.T) {
	repository, err := NewPostgresRepository(pgDriver, stringConnection, zap.L())

	assert.NoError(t, err)
	err = repository.CheckDBConnection()
	assert.NoError(t, err)
}

func TestUsersTable(t *testing.T) {
	repo, err := NewPostgresRepository(pgDriver, stringConnection, zap.L())
	assert.NoError(t, err)

	err = repo.CheckDBConnection()
	require.NoError(t, err)

	err = Reset(repo)
	require.NoError(t, err)

	err = Migrate(repo)
	require.NoError(t, err)

	user := &model.User{
		ID:   1,
		Name: "user",
	}
	userSaved, err := repo.SaveUser(*user)
	assert.NoError(t, err)
	assert.Equal(t, user, userSaved)
}
