package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/MrBorisT/url_shortener/internal/autherr"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	Pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{Pool: pool}
}

func (s *UserStore) RegisterUser(ctx context.Context, userRequest models.UserRequest) error {
	hashedPassword, err := s.hashPassword(userRequest.Password)
	if err != nil {
		return err
	}

	newUser := models.User{
		Email:        userRequest.Email,
		PasswordHash: hashedPassword,
	}
	query := "INSERT INTO users (email, password_hash) VALUES ($1, $2)"
	_, err = s.Pool.Exec(ctx, query, newUser.Email, newUser.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PGCodeUniqueViolation {
			return autherr.ErrUserAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *UserStore) GetUserID(ctx context.Context, userRequest models.UserRequest) (string, error) {
	hashedPassword := ""
	query := "SELECT id, password_hash FROM users WHERE email = $1"
	row := s.Pool.QueryRow(ctx, query, userRequest.Email)

	var userID string
	if err := row.Scan(&userID, &hashedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", autherr.ErrInvalidCredentials
		}
		return "", fmt.Errorf("querying user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userRequest.Password)); err != nil {
		return "", autherr.ErrInvalidCredentials
	}

	return userID, nil
}

func (s *UserStore) hashPassword(password string) (string, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptHash), nil
}
