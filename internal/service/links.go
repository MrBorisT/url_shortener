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

func (s *LinkService) UpdateLink(ctx context.Context, userID string, linkID string, req models.UpdateLinkRequest) (*models.Link, error) {
	originalURL := strings.TrimSpace(req.OriginalURL)
	if originalURL == "" {
		return nil, ErrInvalidOriginalURL
	}

	if linkID == "" {
		return nil, ErrEmptyLinkID
	}

	return s.LinksStore.UpdateLink(ctx, userID, linkID, req)
}

func (s *LinkService) DeleteLink(ctx context.Context, userID string, linkID string) error {
	if linkID == "" {
		return ErrEmptyLinkID
	}
	return s.LinksStore.DeleteLink(ctx, userID, linkID)
}

func (s *LinkService) DisableLink(ctx context.Context, userID string, linkID string) error {
	if linkID == "" {
		return ErrEmptyLinkID
	}
	return s.LinksStore.DisableLink(ctx, userID, linkID)
}

func (s *LinkService) GetLink(ctx context.Context, userID string, linkID string) (*models.Link, error) {
	if linkID == "" {
		return nil, ErrEmptyLinkID
	}
	return s.LinksStore.GetLink(ctx, userID, linkID)
}

func (s *LinkService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	return s.LinksStore.GetOriginalURL(ctx, shortCode)
}

func (s *LinkService) IncrementClickCount(ctx context.Context, shortCode string) error {
	return s.LinksStore.IncrementClickCount(ctx, shortCode)
}

func (s *LinkService) ListLinks(ctx context.Context, userID string) ([]models.Link, error) {
	return s.LinksStore.ListLinks(ctx, userID)
}
