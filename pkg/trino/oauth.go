package trino

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// untilExpirationBuffer is the time before actual expiry at which the token is
// considered "almost expired" and will be refreshed. Must be greater than the
// typical HTTP request timeout (30s).
const untilExpirationBuffer = 60 * time.Second

// token represents an OAuth2 access token with expiration tracking.
type token struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"-"`
}

// isAlmostExpired returns true if the token is missing or will expire within the buffer period.
func (t *token) isAlmostExpired() bool {
	if t.AccessToken == "" {
		return true
	}
	return time.Now().Add(untilExpirationBuffer).After(t.ExpiresAt)
}

// oauthClient handles OAuth2 client credentials flow with per-instance token caching.
// Each datasource instance gets its own oauthClient — no global state.
type oauthClient struct {
	httpClient        *http.Client
	clientID          string
	clientSecret      string
	tokenURL          string
	impersonationUser string

	mu           sync.Mutex
	cachedToken  *token
}

// newOAuthClient creates a new per-instance OAuth client.
func newOAuthClient(httpClient *http.Client, clientID, clientSecret, tokenURL, impersonationUser string) *oauthClient {
	return &oauthClient{
		httpClient:        httpClient,
		clientID:          clientID,
		clientSecret:      clientSecret,
		tokenURL:          tokenURL,
		impersonationUser: impersonationUser,
	}
}

// RoundTrip implements http.RoundTripper. It injects the OAuth bearer token
// and optional impersonation header into every outbound request.
func (c *oauthClient) RoundTrip(req *http.Request) (*http.Response, error) {
	tok, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	if c.impersonationUser != "" {
		req.Header.Set("X-Trino-User", c.impersonationUser)
	}

	return c.httpClient.Transport.RoundTrip(req)
}

// getToken returns a valid access token, refreshing if necessary.
// Uses double-checked locking to minimize lock contention.
func (c *oauthClient) getToken() (*token, error) {
	if c.cachedToken != nil && !c.cachedToken.isAlmostExpired() {
		return c.cachedToken, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cachedToken != nil && !c.cachedToken.isAlmostExpired() {
		return c.cachedToken, nil
	}

	newToken, err := c.retrieveToken()
	if err != nil {
		return nil, err
	}
	c.cachedToken = newToken
	return newToken, nil
}

// retrieveToken fetches a new access token from the OAuth2 token endpoint.
func (c *oauthClient) retrieveToken() (*token, error) {
	log.DefaultLogger.Debug("Retrieving OAuth token", "tokenURL", c.tokenURL)

	values := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"grant_type":    {"client_credentials"},
	}

	resp, err := c.httpClient.PostForm(c.tokenURL, values)
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("token endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	tok := &token{}
	if err := json.Unmarshal(body, tok); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	tok.ExpiresAt = time.Now().Add(time.Second * time.Duration(tok.ExpiresIn))
	log.DefaultLogger.Debug("Token will expire at", "expiresAt", tok.ExpiresAt.Format(time.RFC1123Z))

	return tok, nil
}
