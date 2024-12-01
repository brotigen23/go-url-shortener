package utils

import (
	"bufio"
	"encoding/json"
	"math/rand"
	"os"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

func NewRandomString(size int) string {

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

func LoadLocalAliases(filePath string) ([]model.ShortURL, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(file)

	var aliases []model.ShortURL
	for {
		if data, err := buf.ReadBytes('\n'); err == nil {
			alias := &model.ShortURL{}

			err = json.Unmarshal(data, alias)
			if err != nil {
				return nil, err
			}
			aliases = append(aliases, *alias)
		} else {
			break
		}
	}
	return aliases, nil
}

func Foo(shortRULs []model.ShortURL, users []model.User, userURLs []model.UsersShortURLs, filePath string) error {
	return nil
}

func SaveStorage(aliases []model.ShortURL, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	for _, alias := range aliases {
		data, err := json.Marshal(alias)
		if err != nil {
			return err
		}
		_, err = file.Write(data)
		if err != nil {
			return err
		}

		_, err = file.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
