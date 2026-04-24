package storage

import (
	"context"

	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinksStore struct {
	Pool *pgxpool.Pool
}

func NewLinksStore(pool *pgxpool.Pool) *LinksStore {
	return &LinksStore{pool}
}

func (s *LinksStore) ListLinks(ctx context.Context, userID string) ([]models.Link, error) {
	return nil, nil
}

func (s *LinksStore) GetLink(ctx context.Context, userID string, linkID string) (*models.Link, error) {
	return nil, nil
}

func (s *LinksStore) CreateLink(ctx context.Context, userID string, linkReq models.CreateLinkRequest) (*models.Link, error) {
	return nil, nil
}

func (s *LinksStore) UpdateLink(ctx context.Context, userID string, linkID string, linkReq models.UpdateLinkRequest) (*models.Link, error) {
	return nil, nil
}

func (s *LinksStore) DeleteLink(ctx context.Context, userID string, linkID string) error {
	return nil
}

func (s *LinksStore) DisableLink(ctx context.Context, userID string, linkID string) error {
	return nil
}

func (s *LinksStore) validateLinkID(linkID string) bool {
	_, err := uuid.Parse(linkID)
	return err == nil
}

func (s *LinksStore) generateID() string {
	return uuid.New().String()
}
