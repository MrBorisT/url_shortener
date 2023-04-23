package domain

import (
	"github.com/brianvoe/gofakeit/v6"
)

type URLShortenerService interface {
	GetFullURL(shortURL string) (string, error)
	ShortenURL(fullURL string) (string, error)
}

type URLRepository interface {
	SaveURL(shortURL, fullURL string) error
	GetURL(shortURL string) (string, error)
	ExistsShortURL(fullURL string) *string
}

type service struct {
	URLRepository
	maxShortURLLen int
}

func NewService(repo URLRepository, maxShortURLLen int) *service {
	return &service{
		URLRepository:  repo,
		maxShortURLLen: maxShortURLLen,
	}
}

func (s *service) GetFullURL(shortURL string) (string, error) {
	return s.URLRepository.GetURL(shortURL)
}

func (s *service) ShortenURL(fullURL string) (string, error) {
	if fullURL[:7] != "http://" && fullURL[:8] != "https://" {
		fullURL = "https://" + fullURL
	}
	if shortURL := s.URLRepository.ExistsShortURL(fullURL); shortURL != nil {
		return *shortURL, nil
	}

	shortURL := gofakeit.UUID()[:s.maxShortURLLen]
	if err := s.URLRepository.SaveURL(shortURL, fullURL); err != nil {
		return "", err
	}

	return shortURL, nil
}
