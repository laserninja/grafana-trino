package trino

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestBuildDSN_DefaultUser(t *testing.T) {
	settings := backend.DataSourceInstanceSettings{
		URL: "http://localhost:8080",
	}

	dsn, err := buildDSN(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "http://grafana@localhost:8080"
	if dsn != expected {
		t.Errorf("got %q, want %q", dsn, expected)
	}
}

func TestBuildDSN_BasicAuth(t *testing.T) {
	settings := backend.DataSourceInstanceSettings{
		URL:              "https://trino.example.com:443",
		BasicAuthEnabled: true,
		BasicAuthUser:    "admin",
	}

	dsn, err := buildDSN(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "https://admin@trino.example.com:443"
	if dsn != expected {
		t.Errorf("got %q, want %q", dsn, expected)
	}
}

func TestBuildDSN_InvalidURL(t *testing.T) {
	settings := backend.DataSourceInstanceSettings{
		URL: "://bad-url",
	}

	_, err := buildDSN(settings)
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
