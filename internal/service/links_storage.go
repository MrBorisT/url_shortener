package service

import (
	"context"

	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/google/uuid"
)

type LinksStore interface {
	CreateLink(ctx context.Context, link models.Link) (*models.Link, error)
	GetLink(ctx context.Context, userID uuid.UUID, shortURL string) (*models.Link, error)
	ListLinks(ctx context.Context, userID uuid.UUID) ([]models.Link, error)
	UpdateLink(ctx context.Context, userID uuid.UUID, shortURL string, req models.UpdateLinkRequest) (*models.Link, error)
	DeleteLink(ctx context.Context, userID uuid.UUID, shortURL string) error
	DisableLink(ctx context.Context, userID uuid.UUID, shortURL string) error
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
}
