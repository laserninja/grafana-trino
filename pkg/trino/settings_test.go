package trino

import (
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestLoadSettings_DefaultUser(t *testing.T) {
	config := backend.DataSourceInstanceSettings{
		URL:      "http://localhost:8080",
		JSONData: json.RawMessage(`{}`),
	}

	s, err := loadSettings(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.URL.User.Username() != "grafana" {
		t.Errorf("got user %q, want %q", s.URL.User.Username(), "grafana")
	}
}

func TestLoadSettings_BasicAuth(t *testing.T) {
	config := backend.DataSourceInstanceSettings{
		URL:              "https://trino.example.com:443",
		BasicAuthEnabled: true,
		BasicAuthUser:    "admin",
		JSONData:         json.RawMessage(`{}`),
	}

	s, err := loadSettings(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.URL.User.Username() != "admin" {
		t.Errorf("got user %q, want %q", s.URL.User.Username(), "admin")
	}
}

func TestLoadSettings_InvalidURL(t *testing.T) {
	config := backend.DataSourceInstanceSettings{
		URL:      "://bad-url",
		JSONData: json.RawMessage(`{}`),
	}

	_, err := loadSettings(config)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestLoadSettings_JsonDataFields(t *testing.T) {
	config := backend.DataSourceInstanceSettings{
		URL: "http://localhost:8080",
		JSONData: json.RawMessage(`{
			"enableImpersonation": true,
			"tokenUrl": "https://idp.example.com/token",
			"clientId": "my-client",
			"impersonationUser": "testuser",
			"roles": "system:admin;catalog1:reader",
			"clientTags": "tag1,tag2"
		}`),
		DecryptedSecureJSONData: map[string]string{
			"accessToken":  "secret-token",
			"clientSecret": "secret-client-secret",
		},
	}

	s, err := loadSettings(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !s.EnableImpersonation {
		t.Error("expected enableImpersonation to be true")
	}
	if s.TokenURL != "https://idp.example.com/token" {
		t.Errorf("got tokenUrl %q, want %q", s.TokenURL, "https://idp.example.com/token")
	}
	if s.ClientID != "my-client" {
		t.Errorf("got clientId %q, want %q", s.ClientID, "my-client")
	}
	if s.ImpersonationUser != "testuser" {
		t.Errorf("got impersonationUser %q, want %q", s.ImpersonationUser, "testuser")
	}
	if s.Roles != "system:admin;catalog1:reader" {
		t.Errorf("got roles %q", s.Roles)
	}
	if s.ClientTags != "tag1,tag2" {
		t.Errorf("got clientTags %q", s.ClientTags)
	}
	if s.AccessToken != "secret-token" {
		t.Errorf("got accessToken %q", s.AccessToken)
	}
	if s.ClientSecret != "secret-client-secret" {
		t.Errorf("got clientSecret %q", s.ClientSecret)
	}
}

func TestLoadSettings_InvalidJSON(t *testing.T) {
	config := backend.DataSourceInstanceSettings{
		URL:      "http://localhost:8080",
		JSONData: json.RawMessage(`{invalid}`),
	}

	_, err := loadSettings(config)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
