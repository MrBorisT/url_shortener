package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/shortcode"
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

	rows, err := s.Pool.Query(ctx, query, userID, linkID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (s *LinksStore) CreateLink(ctx context.Context, userID string, linkReq models.CreateLinkRequest) (*models.Link, error) {
	trimmedURL := strings.TrimSpace(linkReq.OriginalURL)
	if trimmedURL == "" {
		return nil, ErrEmptyOriginalURL
	}

	sCode, err := shortcode.Generate()
	if err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	newLink := models.Link{
		OriginalURL: trimmedURL,
		UserID:      userID,
		ShortCode:   sCode,
	}

	query := "INSERT INTO links (original_url, user_id, short_code) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
	row := s.Pool.QueryRow(ctx, query, newLink.OriginalURL, newLink.UserID, newLink.ShortCode)

	if err := row.Scan(&newLink.ID, &newLink.CreatedAt, &newLink.UpdatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PGCodeUniqueViolation {
			//todo retry?
			return nil, fmt.Errorf("creating link: %w", err)
		} else {
			return nil, fmt.Errorf("creating link: %w", err)
		}
	}

	return &newLink, nil
}

func (s *LinksStore) UpdateLink(ctx context.Context, userID string, linkID string, linkReq models.UpdateLinkRequest) (*models.Link, error) {
	//todo validation
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

	_, err := s.Pool.Exec(ctx, query, userID, linkID)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	return nil
}

func (s *LinksStore) DisableLink(ctx context.Context, userID string, linkID string) error {
	query := `
	UPDATE links
	SET disabled_at = NOW()
	WHERE user_id = $1 AND id = $2
	`

	_, err := s.Pool.Exec(ctx, query, userID, linkID)
	if err != nil {
		return fmt.Errorf("disable link: %w", err)
	}

	return nil
}
