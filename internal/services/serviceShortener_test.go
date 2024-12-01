package services

import (
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestServiceShortener(t *testing.T) {

	config := config.NewConfig()
	config.DatabaseDSN = "host=localhost port=5432 user=myuser password=1234 dbname=mydb sslmode=disable"
	service, err := NewService(config, 8, nil, nil, nil)
	assert.NoError(t, err)

	// ---------------------------- Params ----------------------------
	userName := "User1"
	URL := "URL1"

	// ---------------------------- Create User ----------------------------

	user, err := service.repository.SaveUser(*model.NewUser(0, userName))
	assert.NoError(t, err)

	// ---------------------------- Save URL ----------------------------

	alias, err := service.SaveURL(user.Name, URL)
	assert.NoError(t, err)

	// ---------------------------- Get URL ----------------------------

	url, err := service.GetURL(user.Name, alias)
	assert.NoError(t, err)
	assert.Equal(t, url, URL)

	// ---------------------------- Save URL ----------------------------
	URLs := []string{"URL2", "URL3"}

	aliases, err := service.SaveURLs(user.Name, URLs)
	assert.NoError(t, err)

	// ---------------------------- Get URLs ----------------------------
	actualURLs := make(map[string]string)
	actualURLs["URL1"] = alias
	for _, u := range URLs {
		actualURLs[u] = aliases[u]
	}
	urls, err := service.GetURLs(user.Name)
	assert.NoError(t, err)
	assert.Equal(t, urls, actualURLs)
}
