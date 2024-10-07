package storage

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/utils"
)

const (
	ALIASLENGHT = 8
)

type storage struct {
	urls  []string
	alias []string
}

var Storage = NewStorage()

func NewStorage() *storage {
	return &storage{
		urls:  []string{"https://ya.ru", "https://google.com"},
		alias: []string{"asdfgh", "qwerty"},
	}
}

func (s *storage) Put(url []byte) string {
	if s.FindByURL(url) != "" {
		return s.FindByURL(url)
	}
	s.urls = append(s.urls, string(url))
	s.alias = append(s.alias, utils.NewRandomString(ALIASLENGHT))
	return s.alias[len(s.alias)-1]
}

func (s storage) FindByURL(url []byte) string {
	for i, v := range s.urls {
		if string(url) == string(v) {
			return s.alias[i]
		}
	}
	return ""
}
func (s storage) FindByAlias(alias []byte) string {
	for i, v := range s.alias {
		if string(alias) == string(v) {
			return s.urls[i]
		}
	}
	return ""
}

func (s storage) Print() {
	for i := range s.urls {
		fmt.Printf("[%v] %v\n", s.urls[i], s.alias[i])
	}

}

func (s storage) String() string {
	ret := ""
	for i := range s.urls {
		ret += fmt.Sprintf("[%v] %v\n", s.urls[i], s.alias[i])
	}
	return ret
}
