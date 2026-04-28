package service

import (
	"context"

	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/google/uuid"
)

type LinksStore interface {
	CreateLink(ctx context.Context, link models.Link) (*models.Link, error)
	GetLink(ctx context.Context, userID uuid.UUID, linkID string) (*models.Link, error)
	ListLinks(ctx context.Context, userID uuid.UUID) ([]models.Link, error)
	UpdateLink(ctx context.Context, userID uuid.UUID, linkID string, req models.UpdateLinkRequest) (*models.Link, error)
	DeleteLink(ctx context.Context, userID uuid.UUID, linkID string) error
	DisableLink(ctx context.Context, userID uuid.UUID, linkID string) error
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
}
