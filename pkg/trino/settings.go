package trino

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
)

// Settings holds all parsed configuration for a Trino datasource instance.
type Settings struct {
	URL  *url.URL
	Opts httpclient.Options

	// Auth
	EnableImpersonation bool   `json:"enableImpersonation"`
	AccessToken         string `json:"-"` // from secureJsonData
	TokenURL            string `json:"tokenUrl"`
	ClientID            string `json:"clientId"`
	ClientSecret        string `json:"-"` // from secureJsonData
	ImpersonationUser   string `json:"impersonationUser"`
	Roles               string `json:"roles"`
	ClientTags          string `json:"clientTags"`
}

// loadSettings parses the datasource instance settings into a Settings struct.
func loadSettings(config backend.DataSourceInstanceSettings) (*Settings, error) {
	s := &Settings{}

	opts, err := config.HTTPClientOptions(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client options: %w", err)
	}
	s.Opts = opts

	s.URL, err = url.Parse(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %q: %w", config.URL, err)
	}

	// Set user on URL from basic auth or default
	if opts.BasicAuth != nil {
		if opts.BasicAuth.Password != "" {
			s.URL.User = url.UserPassword(opts.BasicAuth.User, opts.BasicAuth.Password)
		} else {
			s.URL.User = url.User(opts.BasicAuth.User)
		}
	} else {
		s.URL.User = url.User("grafana")
	}

	// Parse jsonData fields
	if err := json.Unmarshal(config.JSONData, s); err != nil {
		return nil, fmt.Errorf("failed to parse JSON data: %w", err)
	}

	// Read secrets from secureJsonData
	if token, ok := config.DecryptedSecureJSONData["accessToken"]; ok {
		s.AccessToken = token
	}
	if secret, ok := config.DecryptedSecureJSONData["clientSecret"]; ok {
		s.ClientSecret = secret
	}

	return s, nil
}
