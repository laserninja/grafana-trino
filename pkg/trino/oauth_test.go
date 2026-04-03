package trino

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenIsAlmostExpired_Empty(t *testing.T) {
	tok := &token{}
	if !tok.isAlmostExpired() {
		t.Error("empty token should be almost expired")
	}
}

func TestTokenIsAlmostExpired_Fresh(t *testing.T) {
	tok := &token{
		AccessToken: "test-token",
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
	if tok.isAlmostExpired() {
		t.Error("fresh token should not be almost expired")
	}
}

func TestTokenIsAlmostExpired_ExpiringSoon(t *testing.T) {
	tok := &token{
		AccessToken: "test-token",
		ExpiresAt:   time.Now().Add(30 * time.Second), // Within 60s buffer
	}
	if !tok.isAlmostExpired() {
		t.Error("token expiring within buffer should be almost expired")
	}
}

func TestOAuthClient_RetrieveToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
		}
		if r.FormValue("client_id") != "test-id" {
			t.Errorf("expected client_id=test-id, got %s", r.FormValue("client_id"))
		}
		if r.FormValue("client_secret") != "test-secret" {
			t.Errorf("expected client_secret=test-secret, got %s", r.FormValue("client_secret"))
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"access_token": "test-access-token",
			"expires_in":   3600,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOAuthClient(server.Client(), "test-id", "test-secret", server.URL, "")

	tok, err := client.retrieveToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.AccessToken != "test-access-token" {
		t.Errorf("got access token %q, want %q", tok.AccessToken, "test-access-token")
	}
	if tok.ExpiresIn != 3600 {
		t.Errorf("got expires_in %d, want %d", tok.ExpiresIn, 3600)
	}
}

func TestOAuthClient_GetToken_Cached(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"access_token": "token-" + string(rune('0'+callCount)),
			"expires_in":   3600,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOAuthClient(server.Client(), "id", "secret", server.URL, "")

	// First call should fetch
	tok1, err := client.getToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second call should use cache
	tok2, err := client.getToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 token fetch, got %d", callCount)
	}
	if tok1.AccessToken != tok2.AccessToken {
		t.Error("expected same cached token")
	}
}

func TestOAuthClient_RetrieveToken_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}))
	defer server.Close()

	client := newOAuthClient(server.Client(), "id", "bad-secret", server.URL, "")

	_, err := client.retrieveToken()
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}
