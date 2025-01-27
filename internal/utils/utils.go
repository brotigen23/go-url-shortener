package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/golang-jwt/jwt/v4"
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

func LoadStorage(filePath string) ([]model.ShortURL, error) {
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

type UserJWTClaims struct {
	jwt.RegisteredClaims
	username string
}

func GetUsernameFromJWT(tokenString string, key string) (string, error) {
	claims := &UserJWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	return claims.username, nil
}

func BuildJWTString(username string, key string, expires time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
		},
		username: username,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}
