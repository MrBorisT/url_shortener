package service

import (
	"context"
	"testing"

	"github.com/MrBorisT/url_shortener/internal/linkerr"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/validation"
	"github.com/google/uuid"
)

type fakeLinksStore struct {
	updateLinkFunc     func(ctx context.Context, userID uuid.UUID, shortURL string, req models.UpdateLinkRequest) (*models.Link, error)
	deleteLinkFunc     func(ctx context.Context, userID uuid.UUID, shortURL string) error
	getOriginalURLFunc func(ctx context.Context, shortCode string) (string, error)
}

func (f *fakeLinksStore) UpdateLink(ctx context.Context, userID uuid.UUID, shortURL string, req models.UpdateLinkRequest) (*models.Link, error) {
	return f.updateLinkFunc(ctx, userID, shortURL, req)
}

func (f *fakeLinksStore) CreateLink(ctx context.Context, link models.Link) (*models.Link, error) {
	panic("not implemented")
}

func (f *fakeLinksStore) GetLink(ctx context.Context, userID uuid.UUID, shortURL string) (*models.Link, error) {
	panic("not implemented")
}

func (f *fakeLinksStore) ListLinks(ctx context.Context, userID uuid.UUID) ([]models.Link, error) {
	panic("not implemented")
}

func (f *fakeLinksStore) DeleteLink(ctx context.Context, userID uuid.UUID, shortURL string) error {
	return f.deleteLinkFunc(ctx, userID, shortURL)
}

func (f *fakeLinksStore) DisableLink(ctx context.Context, userID uuid.UUID, shortURL string) error {
	panic("not implemented")
}

func (f *fakeLinksStore) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	return f.getOriginalURLFunc(ctx, shortCode)
}

func (f *fakeLinksStore) IncrementClickCount(ctx context.Context, shortCode string) error {
	panic("not implemented")
}

func TestUpdateLink(t *testing.T) {
	const notFoundLinkURL = "not-found-link-id"
	store := &fakeLinksStore{
		updateLinkFunc: func(ctx context.Context, userID uuid.UUID, shortURL string, req models.UpdateLinkRequest) (*models.Link, error) {
			if shortURL == notFoundLinkURL {
				return nil, linkerr.ErrLinkNotFound
			}
			return &models.Link{
				ID:          uuid.New(),
				UserID:      userID,
				OriginalURL: req.OriginalURL,
			}, nil
		},
	}

	service := NewLinkService(store)
	userID := uuid.New()

	tests := []struct {
		name     string
		shortURL string
		req      models.UpdateLinkRequest
		wantErr  error
		wantURL  string
	}{
		{
			name:     "valid",
			shortURL: "test-url",
			req: models.UpdateLinkRequest{
				OriginalURL: "https://example.com",
			},
			wantErr: nil,
			wantURL: "https://example.com",
		},
		{
			name:     "empty original URL",
			shortURL: "test-url",
			req:      models.UpdateLinkRequest{OriginalURL: ""},
			wantErr:  validation.ErrURLRequired,
			wantURL:  "",
		},
		{
			name:     "invalid original URL",
			shortURL: "test-url",
			req:      models.UpdateLinkRequest{OriginalURL: "not-a-url"},
			wantErr:  validation.ErrURLInvalid,
			wantURL:  "",
		},
		{
			name:     "invalid URL scheme",
			shortURL: "test-url",
			req:      models.UpdateLinkRequest{OriginalURL: "ftp://www.example.com"},
			wantErr:  validation.ErrURLInvalidScheme,
			wantURL:  "",
		},
		{
			name:     "link not found",
			shortURL: notFoundLinkURL,
			req: models.UpdateLinkRequest{
				OriginalURL: "https://example.com",
			},
			wantErr: linkerr.ErrLinkNotFound,
			wantURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.UpdateLink(context.Background(), userID, tt.shortURL, tt.req)

			if tt.wantErr != err {
				t.Fatalf("error check: expected %v got %v", tt.wantErr, err)
			}

			if tt.wantURL != "" && tt.wantURL != got.OriginalURL {
				t.Fatalf("UpdateLink() = %v, want %v", got.OriginalURL, tt.wantURL)
			}
		})
	}
}

func TestDeleteLink(t *testing.T) {
	const notFoundLinkURL = "not-found-link-id"
	store := &fakeLinksStore{
		deleteLinkFunc: func(ctx context.Context, userID uuid.UUID, shortURL string) error {
			if shortURL == notFoundLinkURL {
				return linkerr.ErrLinkNotFound
			}
			return nil
		},
	}

	service := NewLinkService(store)
	userID := uuid.New()

	tests := []struct {
		name     string
		shortURL string
		wantErr  error
	}{
		{
			name:     "valid",
			shortURL: "test-url",
			wantErr:  nil,
		},
		{
			name:     "link not found",
			shortURL: notFoundLinkURL,
			wantErr:  linkerr.ErrLinkNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteLink(context.Background(), userID, tt.shortURL)
			if tt.wantErr != err {
				t.Fatalf("error check: expected %v got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	const notFoundShortCode = "not-found-short-code"
	store := &fakeLinksStore{
		getOriginalURLFunc: func(ctx context.Context, shortCode string) (string, error) {
			if shortCode == notFoundShortCode {
				return "", linkerr.ErrLinkNotFound
			}
			return "https://example.com", nil
		},
	}

	service := NewLinkService(store)
	tests := []struct {
		name      string
		shortCode string
		wantURL   string
		wantErr   error
	}{
		{
			name:      "valid",
			shortCode: "test-short-code",
			wantURL:   "https://example.com",
			wantErr:   nil,
		},
		{
			name:      "short code not found",
			shortCode: notFoundShortCode,
			wantURL:   "",
			wantErr:   linkerr.ErrLinkNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := service.GetOriginalURL(context.Background(), tt.shortCode)
			if tt.wantErr != err {
				t.Fatalf("error check: expected %v got %v", tt.wantErr, err)
			}
			if tt.wantURL != gotURL {
				t.Fatalf("GetOriginalURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}

}
