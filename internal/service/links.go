package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/shortcode"
	"github.com/MrBorisT/url_shortener/internal/storage"
)

const maxShortCodeRetries = 5

type LinkService struct {
	LinksStore *storage.LinksStore
}

func NewLinkService(linksStore *storage.LinksStore) *LinkService {
	return &LinkService{LinksStore: linksStore}
}

func (s *LinkService) CreateLink(ctx context.Context, userID string, req models.CreateLinkRequest) (*models.Link, error) {
	originalURL := strings.TrimSpace(req.OriginalURL)
	if originalURL == "" {
		return nil, ErrInvalidOriginalURL
	}

	for attempt := 0; attempt < maxShortCodeRetries; attempt++ {
		code, err := shortcode.Generate()
		if err != nil {
			return nil, fmt.Errorf("generate short code: %w", err)
		}

		link := models.Link{
			UserID:      userID,
			OriginalURL: originalURL,
			ShortCode:   code,
		}

		createdLink, err := s.LinksStore.CreateLink(ctx, link)
		if err == nil {
			return createdLink, nil
		}

		if errors.Is(err, storage.ErrShortCodeTaken) {
			continue
		}

		return nil, err
	}

	return nil, ErrCouldNotGenerateShortCode
}
