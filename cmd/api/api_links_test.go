//go:build integration

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLinksRoutes(t *testing.T) {
	api := setupTestAPI(t)

	type testUser struct {
		email    string
		password string
	}

	userA := testUser{
		email:    "test@example.com",
		password: "secret123",
	}

	tests := []struct {
		name               string
		httpMethod         string
		url                string
		body               string
		expectedStatus     int
		userB              testUser
		createLinkForUserA bool
	}{
		{
			//user B here is the same as user A
			name:               "create link success",
			httpMethod:         http.MethodPost,
			url:                "/api/links",
			body:               `{"original_url":"https://example.com"}`,
			expectedStatus:     http.StatusCreated,
			userB:              testUser{email: userA.email, password: userA.password},
			createLinkForUserA: false,
		},
		{
			name:               "user B cannot access user A's links",
			httpMethod:         http.MethodGet,
			url:                "/api/links/{link_id}",
			body:               "",
			expectedStatus:     http.StatusNotFound,
			userB:              testUser{email: "test-1@example.com", password: "secret123"},
			createLinkForUserA: true,
		},
		{
			name:               "user B cannot update user A's links",
			httpMethod:         http.MethodPatch,
			url:                "/api/links/{link_id}",
			body:               `{"original_url":"https://example.com"}`,
			expectedStatus:     http.StatusNotFound,
			userB:              testUser{email: "test-1@example.com", password: "secret123"},
			createLinkForUserA: true,
		},
		{
			name:               "user B cannot delete user A's links",
			httpMethod:         http.MethodDelete,
			url:                "/api/links/{link_id}",
			body:               "",
			expectedStatus:     http.StatusNotFound,
			userB:              testUser{email: "test-1@example.com", password: "secret123"},
			createLinkForUserA: true,
		},
		{
			name:               "user B cannot disable user A's links",
			httpMethod:         http.MethodPost,
			url:                "/api/links/{link_id}/disable",
			body:               "",
			expectedStatus:     http.StatusNotFound,
			userB:              testUser{email: "test-1@example.com", password: "secret123"},
			createLinkForUserA: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			truncateTestDB(t, api.pool)

			registerUser(t, api.router, tt.userB.email, tt.userB.password)
			tokenUserB := loginUser(t, api.router, tt.userB.email, tt.userB.password)

			if tt.createLinkForUserA {
				registerUser(t, api.router, userA.email, userA.password)
				tokenUserA := loginUser(t, api.router, userA.email, userA.password)
				link := createTestLink(t, api.router, tokenUserA, "https://example.com")
				tt.url = strings.Replace(tt.url, "{link_id}", link.ID.String(), 1)
			}

			req := httptest.NewRequest(tt.httpMethod, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tokenUserB)

			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d, body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}
