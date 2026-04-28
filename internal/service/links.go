package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/MrBorisT/url_shortener/internal/linkerr"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/shortcode"
	"github.com/MrBorisT/url_shortener/internal/validation"
	"github.com/google/uuid"
)

const maxShortCodeRetries = 5

type LinkService struct {
	LinksStore LinksStore
}

func NewLinkService(linksStore LinksStore) *LinkService {
	return &LinkService{LinksStore: linksStore}
}

func (s *LinkService) CreateLink(ctx context.Context, userID uuid.UUID, req models.CreateLinkRequest) (*models.Link, error) {
	originalURL, err := validation.NormalizeURL(req.OriginalURL)
	if err != nil {
		return nil, err
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

		if errors.Is(err, linkerr.ErrShortCodeTaken) {
			continue
		}

		return nil, err
	}

	return nil, ErrCouldNotGenerateShortCode
}

func (s *LinkService) UpdateLink(ctx context.Context, userID uuid.UUID, linkID string, req models.UpdateLinkRequest) (*models.Link, error) {
	if linkID == "" {
		return nil, ErrEmptyLinkID
	}

	originalURL, err := validation.NormalizeURL(req.OriginalURL)
	if err != nil {
		return nil, err
	}

	req.OriginalURL = originalURL

	link, err := s.LinksStore.UpdateLink(ctx, userID, linkID, req)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}
	return link, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, userID uuid.UUID, linkID string) error {
	if linkID == "" {
		return ErrEmptyLinkID
	}

	if err := s.LinksStore.DeleteLink(ctx, userID, linkID); err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return ErrLinkNotFound
		}
		return err
	}
	return nil
}

func (s *LinkService) DisableLink(ctx context.Context, userID uuid.UUID, linkID string) error {
	if linkID == "" {
		return ErrEmptyLinkID
	}
	if err := s.LinksStore.DisableLink(ctx, userID, linkID); err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return ErrLinkNotFound
		}
		return err
	}
	return nil
}

func (s *LinkService) GetLink(ctx context.Context, userID uuid.UUID, linkID string) (*models.Link, error) {
	if linkID == "" {
		return nil, ErrEmptyLinkID
	}
	link, err := s.LinksStore.GetLink(ctx, userID, linkID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}
	return link, nil
}

func (s *LinkService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	return s.LinksStore.GetOriginalURL(ctx, shortCode)
}

func (s *LinkService) IncrementClickCount(ctx context.Context, shortCode string) error {
	return s.LinksStore.IncrementClickCount(ctx, shortCode)
}

func (s *LinkService) ListLinks(ctx context.Context, userID uuid.UUID) ([]models.Link, error) {
	return s.LinksStore.ListLinks(ctx, userID)
}
