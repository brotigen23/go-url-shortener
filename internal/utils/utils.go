package utils

import (
	"bufio"
	"encoding/json"
	"math/rand"
	"os"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

const FILENAME = "aliases.txt"

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

func LoadLocalAliases(filePath string) ([]model.Alias, error) {
	file, _ := os.OpenFile(filePath+FILENAME, os.O_RDONLY|os.O_CREATE, 0666)
	buf := bufio.NewReader(file)

	var aliases []model.Alias
	for {
		if data, err := buf.ReadBytes('\n'); err == nil {
			alias := &model.Alias{}

			json.Unmarshal(data, alias)
			aliases = append(aliases, *alias)
		} else {
			break
		}
	}
	return aliases, nil
}

func SaveLocalAliases(aliases []model.Alias, filePath string) error {
	return nil
	file, _ := os.OpenFile(filePath+FILENAME, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	for _, alias := range aliases {
		data, _ := json.Marshal(alias)
		file.Write(data)
		file.Write([]byte("\n"))
	}
	return nil
}
