package service

import (
	"context"
	"net/mail"
	"strings"

	auth "github.com/MrBorisT/url_shortener/internal/jwt"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/storage"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	AuthStore  *storage.UserStore
	JWTManager *auth.JWTManager
}

func NewAuthService(AuthStore *storage.UserStore, JWTManager *auth.JWTManager) *AuthService {
	return &AuthService{AuthStore: AuthStore, JWTManager: JWTManager}
}

func (s *AuthService) RegisterUser(ctx context.Context, userRequest models.UserRequest) error {
	if err := verifyUserRequest(&userRequest); err != nil {
		return err
	}
	if err := s.AuthStore.RegisterUser(ctx, userRequest); err != nil {
		if err == storage.ErrUserAlreadyExists {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (s *AuthService) AuthenticateUser(ctx context.Context, userRequest models.UserRequest) (string, error) {
	userID, err := s.AuthStore.GetUserID(ctx, userRequest)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *AuthService) GenerateJWT(userID string) (string, error) {
	return s.JWTManager.GenerateJWT(userID)
}

func verifyUserRequest(userRequest *models.UserRequest) error {
	trimmedEmail := strings.TrimSpace(userRequest.Email)
	if trimmedEmail == "" {
		return ErrEmptyUserEmail
	}
	if _, err := mail.ParseAddress(trimmedEmail); err != nil {
		return ErrIncorrectUserEmail
	}

	trimmedPassword := strings.TrimSpace(userRequest.Password)
	if trimmedPassword == "" {
		return ErrEmptyUserPassword
	}
	if len(trimmedPassword) < 6 {
		return ErrShortUserPassword
	}
	if len(trimmedPassword) > 72 {
		return ErrLongUserPassword
	}

	userRequest.Email = trimmedEmail
	userRequest.Password = trimmedPassword

	return nil
}

func (s *AuthService) Verify(tokenString string) (*jwt.RegisteredClaims, error) {
	return s.JWTManager.Verify(tokenString)
}
