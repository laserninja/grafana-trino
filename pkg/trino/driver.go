package trino

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/trinodb/trino-go-client/trino"
)

// openDB creates a sql.DB connection to Trino using the parsed settings.
// It configures TLS, OAuth2 (if needed), roles, and the Trino client.
func openDB(settings *Settings) (*sql.DB, error) {
	transport, err := buildTransport(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to build transport: %w", err)
	}

	httpClient := &http.Client{Transport: transport}

	// If OAuth2 settings are configured, wrap the HTTP client with the OAuth transport.
	if settings.TokenURL != "" || settings.ClientID != "" || settings.ClientSecret != "" {
		if settings.AccessToken != "" {
			return nil, fmt.Errorf("access token must not be set when using OAuth2 authentication")
		}

		var missingParams []string
		if settings.TokenURL == "" {
			missingParams = append(missingParams, "Token URL")
		}
		if settings.ClientID == "" {
			missingParams = append(missingParams, "Client ID")
		}
		if settings.ClientSecret == "" {
			missingParams = append(missingParams, "Client secret")
		}
		if len(missingParams) > 0 {
			return nil, fmt.Errorf("missing OAuth2 parameters: %s", strings.Join(missingParams, ", "))
		}

		oauthTransport := newOAuthClient(httpClient, settings.ClientID, settings.ClientSecret, settings.TokenURL, settings.ImpersonationUser)
		httpClient = &http.Client{Transport: oauthTransport}
	}

	if err := trino.RegisterCustomClient("grafana", httpClient); err != nil {
		return nil, fmt.Errorf("failed to register Trino HTTP client: %w", err)
	}

	roles, err := parseRoles(settings.Roles)
	if err != nil {
		return nil, err
	}

	config := trino.Config{
		ServerURI:                  settings.URL.String(),
		Source:                     "grafana",
		CustomClientName:           "grafana",
		ForwardAuthorizationHeader: true,
		AccessToken:                settings.AccessToken,
		Roles:                      roles,
	}

	dsn, err := config.FormatDSN()
	if err != nil {
		return nil, fmt.Errorf("failed to format Trino DSN: %w", err)
	}

	return sql.Open("trino", dsn)
}

// buildTransport creates an http.Transport with TLS configuration from settings.
func buildTransport(settings *Settings) (*http.Transport, error) {
	skipVerify := false
	var certPool *x509.CertPool
	var clientCerts []tls.Certificate

	if settings.Opts.TLS != nil {
		skipVerify = settings.Opts.TLS.InsecureSkipVerify

		if settings.Opts.TLS.CACertificate != "" {
			certPool = x509.NewCertPool()
			if !certPool.AppendCertsFromPEM([]byte(settings.Opts.TLS.CACertificate)) {
				return nil, fmt.Errorf("failed to parse CA certificate")
			}
		}

		if settings.Opts.TLS.ClientCertificate != "" {
			if settings.Opts.TLS.ClientKey == "" {
				return nil, fmt.Errorf("client certificate was configured without a client key")
			}
			cert, err := tls.X509KeyPair(
				[]byte(settings.Opts.TLS.ClientCertificate),
				[]byte(settings.Opts.TLS.ClientKey),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to load client certificate: %w", err)
			}
			clientCerts = append(clientCerts, cert)
		}
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipVerify,
			Certificates:       clientCerts,
			RootCAs:            certPool,
		},
	}, nil
}

// parseRoles parses the "catalog:role;catalog2:role2" format into a map.
func parseRoles(roleStr string) (map[string]string, error) {
	roles := make(map[string]string)
	if strings.TrimSpace(roleStr) == "" {
		return roles, nil
	}

	pairs := strings.Split(roleStr, ";")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid role format: expected catalog:role, got %q", pair)
		}
		catalog := strings.TrimSpace(parts[0])
		role := strings.TrimSpace(parts[1])
		if catalog != "" && role != "" {
			roles[catalog] = role
		}
	}

	return roles, nil
}
