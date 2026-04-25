package models

import (
	"time"

	"github.com/google/uuid"
)

// db link
type Link struct {
	ID          uuid.UUID  `json:"id"`
	UserID      string     `json:"user_id"`
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code"`
	ClickCount  int64      `json:"click_count"`
	DisabledAt  *time.Time `json:"disabled_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// requests
type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

// maybe add more functionality later
type UpdateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}
