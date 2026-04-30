package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToken(t *testing.T) {
	api := setupTestAPI(t)

	truncateTestDB(t, api.pool)
	registerUser(t, api.router, "test-1@example.com", "secret123")
	token := loginUser(t, api.router, "test-1@example.com", "secret123")

	tests := []struct {
		name       string
		httpMethod string
		wantCode   int
		token      string
	}{
		{
			name:       "protected route without token",
			httpMethod: http.MethodGet,
			wantCode:   http.StatusUnauthorized,
			token:      "",
		},
		{
			name:       "protected route malformed token",
			httpMethod: http.MethodGet,
			wantCode:   http.StatusUnauthorized,
			token:      token + "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.httpMethod, "/api/links", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)

			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Fatalf("expected status %d, got %d, body: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
		})
	}
}
