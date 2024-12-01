package repositories

import (
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	var repository Repository
	repository, err := NewPostgresRepository("postgres", "host=localhost port=5432 user=myuser password=1234 dbname=mydb sslmode=disable", nil)
	assert.NoError(t, err)
	// ------------------------------ USERS ------------------------------
	usersToSave := []model.User{
		*model.NewUser(0, "User1"),
		*model.NewUser(0, "User2"),
		*model.NewUser(0, "User3"),
		*model.NewUser(0, "User4"),
	}
	var users []*model.User
	// Save users
	for _, user := range usersToSave {
		u, err := repository.SaveUser(user)
		assert.NoError(t, err)
		users = append(users, u)
	}
	// Get users
	for _, user := range users {
		_, err := repository.GetUserByName(user.Name)
		assert.NoError(t, err)
	}

	// ------------------------------ SHORT URLs ------------------------------
	shortURLsToSave := []model.ShortURL{
		*model.NewShortURL(0, "URL1", "Alias1"),
		*model.NewShortURL(0, "URL2", "Alias2"),
		*model.NewShortURL(0, "URL3", "Alias3"),
		*model.NewShortURL(0, "URL4", "Alias4"),
	}
	var shortURLs []*model.ShortURL
	// Save URLs
	for _, url := range shortURLsToSave {
		shortURL, err := repository.SaveShortURL(url)
		assert.NoError(t, err)

		shortURLs = append(shortURLs, shortURL)
	}
	// Get URLs
	for _, shortURL := range shortURLs {
		_, err := repository.GetShortURLByURL(shortURL.URL)
		assert.NoError(t, err)
	}

	// ------------------------------ USERS_URL ------------------------------
	for i := range users {
		userID := users[i].Id
		shortURLID := shortURLs[i].Id
		userShortURL := model.NewUsers_ShortURLs(0, userID, shortURLID)
		_, err = repository.SaveUserShortURL(*userShortURL)
		assert.NoError(t, err)
	}
}
