package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinksStore struct {
	Pool *pgxpool.Pool
}

func NewLinksStore(pool *pgxpool.Pool) *LinksStore {
	return &LinksStore{pool}
}

func (s *LinksStore) ListLinks(ctx context.Context, userID string) ([]models.Link, error) {
	query := `
	SELECT id, original_url, short_code, click_count, disabled_at, created_at, updated_at FROM links
	WHERE user_id = $1
	`

	rows, err := s.Pool.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultLinks := make([]models.Link, 0)
	for rows.Next() {
		var link models.Link
		if err := rows.Scan(
			&link.ID,
			&link.OriginalURL,
			&link.ShortCode,
			&link.ClickCount,
			&link.DisabledAt,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, err
		}
		resultLinks = append(resultLinks, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultLinks, nil
}

func (s *LinksStore) GetLink(ctx context.Context, userID string, linkID string) (*models.Link, error) {
	query := `
	SELECT id, original_url, short_code, click_count, disabled_at, created_at, updated_at FROM links
	WHERE user_id = $1 AND id = $2
	`

	link := models.Link{}

	if err := s.Pool.QueryRow(ctx, query, userID, linkID).Scan(
		&link.ID,
		&link.OriginalURL,
		&link.ShortCode,
		&link.ClickCount,
		&link.DisabledAt,
		&link.CreatedAt,
		&link.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrLinkNotFound
		} else {
			return nil, fmt.Errorf("get link: %w", err)
		}
	}

	return &link, nil
}

func (s *LinksStore) CreateLink(ctx context.Context, link models.Link) (*models.Link, error) {
	query := `
		INSERT INTO links (original_url, user_id, short_code)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	row := s.Pool.QueryRow(ctx, query, link.OriginalURL, link.UserID, link.ShortCode)

	if err := row.Scan(&link.ID, &link.CreatedAt, &link.UpdatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == PGCodeUniqueViolation &&
			pgErr.ConstraintName == "links_short_code_key" {
			return nil, ErrShortCodeTaken
		}

		return nil, fmt.Errorf("create link: %w", err)
	}

	return &link, nil
}

func (s *LinksStore) UpdateLink(ctx context.Context, userID string, linkID string, linkReq models.UpdateLinkRequest) (*models.Link, error) {
	query := `
	UPDATE links
	SET original_url = $1, updated_at = NOW()
	WHERE user_id = $2 AND id = $3
	RETURNING id, original_url, short_code, click_count, disabled_at, created_at, updated_at
	`

	updatedLink := models.Link{}
	if err := s.Pool.QueryRow(ctx, query, linkReq.OriginalURL, userID, linkID).Scan(
		&updatedLink.ID,
		&updatedLink.OriginalURL,
		&updatedLink.ShortCode,
		&updatedLink.ClickCount,
		&updatedLink.DisabledAt,
		&updatedLink.CreatedAt,
		&updatedLink.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrLinkNotFound
		} else {
			return nil, fmt.Errorf("update link: %w", err)
		}
	}

	return &updatedLink, nil
}

func (s *LinksStore) DeleteLink(ctx context.Context, userID string, linkID string) error {
	query := `
	DELETE FROM links
	WHERE user_id = $1 AND id = $2
	`

	tag, err := s.Pool.Exec(ctx, query, userID, linkID)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrLinkNotFound
	}

	return nil
}

func (s *LinksStore) DisableLink(ctx context.Context, userID string, linkID string) error {
	query := `
	UPDATE links
	SET disabled_at = NOW()
	WHERE user_id = $1 AND id = $2
	`

	tag, err := s.Pool.Exec(ctx, query, userID, linkID)
	if err != nil {
		return fmt.Errorf("disable link: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrLinkNotFound
	}

	return nil
}

func (s *LinksStore) GetOriginalURL(ctx context.Context, shortLink string) (string, error) {
	query := `
	SELECT original_url, disabled_at FROM links
	WHERE short_code = $1
	`

	var originalURL string
	var disabledAt *time.Time
	if err := s.Pool.QueryRow(ctx, query, shortLink).Scan(&originalURL, &disabledAt); err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrLinkNotFound
		} else {
			return "", fmt.Errorf("get original URL: %w", err)
		}
	}

	if disabledAt != nil {
		return "", ErrLinkDisabled
	}

	return originalURL, nil
}

func (s *LinksStore) IncrementClickCount(ctx context.Context, shortLink string) error {
	query := `
	UPDATE links
	SET click_count = click_count + 1, updated_at = NOW()
	WHERE short_code = $1
	`

	_, err := s.Pool.Exec(ctx, query, shortLink)
	if err != nil {
		return fmt.Errorf("increment click count: %w", err)
	}

	return nil
}
