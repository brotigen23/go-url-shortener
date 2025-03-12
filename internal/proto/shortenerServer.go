package proto

import (
	context "context"
	"strconv"

	"github.com/brotigen23/go-url-shortener/internal/service"
)

type ShortenerServer struct {
	UnimplementedShotenerServer

	service *service.Service
}

func NewShortenerServer(service *service.Service) *ShortenerServer {
	return &ShortenerServer{
		service: service,
	}
}

func (s ShortenerServer) CreateShortURL(ctx context.Context, r *SaveURLRequest) (*SaveURLResponse, error) {
	shortURL, err := s.service.CreateShortURL(r.Username, r.Url)
	ret := &SaveURLResponse{
		Username: r.Username,
		Url:      shortURL,
	}
	return ret, err
}

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
