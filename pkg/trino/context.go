package trino

import (
	"context"
	"database/sql"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// Context keys for passing Trino-specific values through the request context.
type contextKey string

const (
	ctxKeyTrinoUser  contextKey = "X-Trino-User"
	ctxKeyAccessToken contextKey = "accessToken"
	ctxKeyClientTags contextKey = "X-Trino-Client-Tags"

	bearerPrefix = "Bearer "
)

// MutateQueryData implements sqlds.QueryDataMutator. It is called by sqlds
// before query execution to inject Trino-specific context values.
func (d *Datasource) MutateQueryData(ctx context.Context, req *backend.QueryDataRequest) (context.Context, *backend.QueryDataRequest) {
	config := req.PluginContext.DataSourceInstanceSettings
	if config == nil {
		return ctx, req
	}

	settings, err := loadSettings(*config)
	if err != nil {
		log.DefaultLogger.Error("Failed to load settings in MutateQueryData", "error", err)
		return ctx, req
	}

	// Inject OAuth token from forwarded headers
	ctx = injectAccessToken(ctx, req)

	// Inject impersonated user
	if settings.EnableImpersonation {
		user := req.PluginContext.User
		if user == nil {
			log.DefaultLogger.Error("User is nil but impersonation is enabled")
		} else {
			ctx = context.WithValue(ctx, ctxKeyTrinoUser, user)
		}
	}

	// Inject client tags
	if settings.ClientTags != "" {
		ctx = context.WithValue(ctx, ctxKeyClientTags, settings.ClientTags)
	}

	return ctx, req
}

// injectAccessToken extracts the OAuth identity token from forwarded headers
// and injects it into the context for use by SetQueryArgs.
func injectAccessToken(ctx context.Context, req *backend.QueryDataRequest) context.Context {
	header := req.GetHTTPHeader(backend.OAuthIdentityTokenHeaderName)
	if strings.HasPrefix(header, bearerPrefix) {
		tok := strings.TrimPrefix(header, bearerPrefix)
		return context.WithValue(ctx, ctxKeyAccessToken, tok)
	}
	return ctx
}

// setQueryArgs reads Trino-specific values from the context and returns them
// as named SQL args for the trino-go-client driver.
func setQueryArgs(ctx context.Context) []interface{} {
	var args []interface{}

	if user := ctx.Value(ctxKeyTrinoUser); user != nil {
		if u, ok := user.(*backend.User); ok {
			log.DefaultLogger.Debug("Setting Trino user from impersonation", "user", u.Login)
			args = append(args, sql.Named(string(ctxKeyTrinoUser), u.Login))
		}
	}

	if tok := ctx.Value(ctxKeyAccessToken); tok != nil {
		if t, ok := tok.(string); ok {
			args = append(args, sql.Named(string(ctxKeyAccessToken), t))
		}
	}

	if tags := ctx.Value(ctxKeyClientTags); tags != nil {
		if t, ok := tags.(string); ok {
			args = append(args, sql.Named(string(ctxKeyClientTags), t))
		}
	}

	return args
}

// Compile-time check that Datasource satisfies QueryDataMutator.
var _ interface {
	MutateQueryData(context.Context, *backend.QueryDataRequest) (context.Context, *backend.QueryDataRequest)
} = (*Datasource)(nil)
