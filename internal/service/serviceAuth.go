package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/utils"
)

type ServiceAuth struct {
	repo repository.Repository

	key    []byte
	signes map[string][]byte
}

func generateKey(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewServiceAuth(config *config.Config, repository repository.Repository) (*ServiceAuth, error) {
	key, err := generateKey(16)
	if err != nil {
		return nil, err
	}
	return &ServiceAuth{
		signes: make(map[string][]byte),
		repo:   repository,
		key:    key,
	}, nil
}

func (s *ServiceAuth) SaveUser(userName string) error {
	_, err := s.repo.SaveUser(*model.NewUser(0, userName))
	return err
}

func (s *ServiceAuth) SignUser(userID string) error {
	h := hmac.New(sha256.New, s.key)
	_, err := h.Write([]byte(userID))
	if err != nil {
		return err
	}
	s.signes[userID] = h.Sum(nil)
	return nil
}

func (s *ServiceAuth) CheckSing(userID string) bool {
	h := hmac.New(sha256.New, s.key)
	h.Write([]byte(userID))
	sign := h.Sum(nil)
	return hmac.Equal(s.signes[userID], sign)
}

func (s *ServiceAuth) GenerateID() (string, error) {

	return utils.NewRandomString(16), nil
}
