package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/MrBorisT/url_shortener/internal/config"
	auth "github.com/MrBorisT/url_shortener/internal/jwt"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/service"
	"github.com/MrBorisT/url_shortener/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type testAPI struct {
	router http.Handler
	pool   *pgxpool.Pool
}

func setupTestAPI(t *testing.T) *testAPI {
	t.Helper()

	_ = godotenv.Load("./.env.test")

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	userStore := storage.NewUserStore(pool)
	cfg := &config.Config{
		JWTSecret: "test-secret",
		JWTTTL:    5 * time.Minute,
	}
	authManager := auth.NewJWTManager(cfg)
	linksStore := storage.NewPostgresLinksStore(pool)
	linkService := service.NewLinkService(linksStore)
	authService := service.NewAuthService(userStore, authManager)

	return &testAPI{
		router: newRouter(linkService, authService),
		pool:   pool,
	}
}

func truncateTestDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(context.Background(), `
		TRUNCATE users, links RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("truncate db: %v", err)
	}
}

func registerUser(t *testing.T, router http.Handler, email, password string) {
	t.Helper()

	body := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("seed register user: expected status %d, got %d, body: %s",
			http.StatusCreated,
			rr.Code,
			rr.Body.String(),
		)
	}
}

func loginUser(t *testing.T, router http.Handler, email, password string) string {
	t.Helper()

	body := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("login user: expected status %d, got %d, body: %s",
			http.StatusOK,
			rr.Code,
			rr.Body.String(),
		)
	}

	var resp struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode login response: %v, body: %s", err, rr.Body.String())
	}

	if resp.Token == "" {
		t.Fatalf("login response has empty token, body: %s", rr.Body.String())
	}

	return resp.Token
}

func createTestLink(t *testing.T, router http.Handler, token, originalURL string) models.Link {
	t.Helper()
	body := fmt.Sprintf(`{"original_url":%q}`, originalURL)

	req := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("create test link: expected status %d, got %d, body: %s",
			http.StatusCreated,
			rr.Code,
			rr.Body.String(),
		)
	}

	var resp models.Link

	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode create link response: %v, body: %s", err, rr.Body.String())
	}

	return resp
}
