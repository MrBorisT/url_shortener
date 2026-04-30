package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrBorisT/url_shortener/internal/models"
)

func TestRedirectSuccess(t *testing.T) {
	api := setupTestAPI(t)
	truncateTestDB(t, api.pool)

	//setting up test data
	originalURL := "https://example.com"
	userEmail := "test@example.com"
	userPassword := "secret123"
	registerUser(t, api.router, userEmail, userPassword)
	token := loginUser(t, api.router, userEmail, userPassword)
	link := createTestLink(t, api.router, token, originalURL)

	req := httptest.NewRequest(http.MethodGet, "/"+link.ShortCode, nil)
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusFound {
		t.Fatalf("expected status %d, got %d, body: %s",
			http.StatusFound,
			rr.Code,
			rr.Body.String(),
		)
	}
}

func TestRedirectDisabledLink(t *testing.T) {
	api := setupTestAPI(t)
	truncateTestDB(t, api.pool)

	//setting up test data
	originalURL := "https://example.com"
	userEmail := "test@example.com"
	userPassword := "secret123"
	registerUser(t, api.router, userEmail, userPassword)
	token := loginUser(t, api.router, userEmail, userPassword)
	link := createTestLink(t, api.router, token, originalURL)
	disableTestLink(t, api.router, token, link.ID.String())

	req := httptest.NewRequest(http.MethodGet, "/"+link.ShortCode, nil)
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusGone {
		t.Fatalf("expected status %d, got %d, body: %s",
			http.StatusGone,
			rr.Code,
			rr.Body.String(),
		)
	}
}

func TestRedirectUnknownLink(t *testing.T) {
	api := setupTestAPI(t)
	truncateTestDB(t, api.pool)

	req := httptest.NewRequest(http.MethodGet, "/"+"unknown", nil)
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d, body: %s",
			http.StatusNotFound,
			rr.Code,
			rr.Body.String(),
		)
	}
}

func TestRedirectIncrements(t *testing.T) {
	api := setupTestAPI(t)
	truncateTestDB(t, api.pool)

	//setting up test data
	originalURL := "https://example.com"
	userEmail := "test@example.com"
	userPassword := "secret123"
	registerUser(t, api.router, userEmail, userPassword)
	token := loginUser(t, api.router, userEmail, userPassword)
	link := createTestLink(t, api.router, token, originalURL)

	clickCount := 67
	for i := 0; i < clickCount; i++ {
		req := httptest.NewRequest(http.MethodGet, "/"+link.ShortCode, nil)
		rr := httptest.NewRecorder()
		api.router.ServeHTTP(rr, req)
		if rr.Code != http.StatusFound {
			t.Fatalf("expected status %d, got %d, body: %s",
				http.StatusFound,
				rr.Code,
				rr.Body.String(),
			)
		}
	}

	//fetch the link details to check redirect count
	req := httptest.NewRequest(http.MethodGet, "/api/links/"+link.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("fetch link details: expected status %d, got %d, body: %s",
			http.StatusOK,
			rr.Code,
			rr.Body.String(),
		)
	}
	var resp models.Link
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode link details response: %v, body: %s", err, rr.Body.String())
	}

	if resp.ClickCount != int64(clickCount) {
		t.Fatalf("expected click count %d, got %d", clickCount, resp.ClickCount)
	}
}
