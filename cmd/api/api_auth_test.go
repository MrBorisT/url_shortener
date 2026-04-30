package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	api := setupTestAPI(t)

	tests := []struct {
		name           string
		httpMethod     string
		url            string
		body           string
		expectedStatus int
	}{
		{
			name:       "valid registration",
			httpMethod: http.MethodPost,
			url:        "/api/auth/register",
			body: `{
				"email": "test-success@example.com",
				"password": "secret123"
			}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:       "duplicate email",
			httpMethod: http.MethodPost,
			url:        "/api/auth/register",
			body: `{
				"email": "test@example.com",
				"password": "secret123"
			}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:       "login success",
			httpMethod: http.MethodPost,
			url:        "/api/auth/login",
			body: `{
				"email": "test@example.com",
				"password": "secret123"
			}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:       "login wrong password",
			httpMethod: http.MethodPost,
			url:        "/api/auth/login",
			body: `{
				"email": "test@example.com",
				"password": "wrongpassword"
			}`,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			truncateTestDB(t, api.pool)
			registerUser(t, api.router, "test@example.com", "secret123")

			req := httptest.NewRequest(tt.httpMethod, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d, body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}