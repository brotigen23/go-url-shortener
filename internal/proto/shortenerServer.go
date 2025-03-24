package proto

import (
	context "context"
	"strconv"

	"github.com/brotigen23/go-url-shortener/internal/service"
)

// Proto server struct
type ShortenerServer struct {
	UnimplementedShotenerServer

	service *service.Service
}

// Return proto server
func NewShortenerServer(service *service.Service) *ShortenerServer {
	return &ShortenerServer{
		service: service,
	}
}

// Create short URL
func (s ShortenerServer) CreateShortURL(ctx context.Context, r *SaveURLRequest) (*SaveURLResponse, error) {
	shortURL, err := s.service.CreateShortURL(r.Username, r.Url)
	ret := &SaveURLResponse{
		Username: r.Username,
		Url:      shortURL,
	}
	return ret, err
}

// Create short URLs
func (s ShortenerServer) CreateShortURLs(ctx context.Context, r *BatchURLRequest) (*BatchURLResponse, error) {
	shortURLs, err := s.service.CreateShortURLs(r.Username, r.Urls)
	if err != nil {
		return nil, err
	}
	mas := make([]string, len(shortURLs))
	for k, v := range shortURLs {
		i, er := strconv.Atoi(k)
		if er != nil {
			return nil, er
		}
		mas[i] = v
	}
	ret := &BatchURLResponse{
		Username: r.Username,
		Urls:     mas,
	}
	return ret, err
}

// Returns short URLs
func (s ShortenerServer) GetShortURL(ctx context.Context, r *GetShortURLRequest) (*GetShortURLResponse, error) {
	shortURL, err := s.service.GetShortURL(r.URL)
	if err != nil {
		return nil, err
	}

	ret := &GetShortURLResponse{
		ShortURL: shortURL,
	}

	return ret, nil
}
