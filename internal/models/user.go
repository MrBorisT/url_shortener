package models

import (
	"time"

	"github.com/google/uuid"
)

// db user
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// requests
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTToken struct {
	Token string `json:"token,omitempty"`
}
