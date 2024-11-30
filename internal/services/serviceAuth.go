package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/repositories"
	"github.com/brotigen23/go-url-shortener/internal/utils"
)

type ServiceAuth struct {
	repo repositories.Repository

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

func NewServiceAuth(config *config.Config) (*ServiceAuth, error) {
	key, err := generateKey(16)
	if err != nil {
		return nil, err
	}
	repo, err := repositories.NewPostgresRepository("postgres", config.DatabaseDSN, nil)
	if err != nil {
		return nil, err
	}

	return &ServiceAuth{
		signes: make(map[string][]byte),
		repo:   repo,
		key:    key,
	}, nil
}

func (s *ServiceAuth) SaveUser(id string) error {
	return s.repo.SaveUser(id)
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

func (s *ServiceAuth) IsExist(userID string) error {
	return s.repo.GetUserByID(userID)
}

func (s *ServiceAuth) GenerateID() (string, error) {

	return utils.NewRandomString(16), nil
}
