package trino

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/sqlds/v4"
	_ "github.com/trinodb/trino-go-client/trino"
)

// Datasource implements the sqlds.Driver interface for Trino.
type Datasource struct{}

// Settings returns the driver settings for Trino.
func (d *Datasource) Settings(_ context.Context, _ backend.DataSourceInstanceSettings) sqlds.DriverSettings {
	return sqlds.DriverSettings{
		FillMode: &fillModeNull,
	}
}

var fillModeNull = data.FillMissing{Mode: data.FillModeNull}

// Connect opens a SQL connection to Trino using the datasource settings.
func (d *Datasource) Connect(_ context.Context, settings backend.DataSourceInstanceSettings, _ json.RawMessage) (*sql.DB, error) {
	dsn, err := buildDSN(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to build Trino DSN: %w", err)
	}

	db, err := sql.Open("trino", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open Trino connection: %w", err)
	}

	return db, nil
}

// Converters returns the type converters for Trino SQL types.
func (d *Datasource) Converters() []sqlutil.Converter {
	return []sqlutil.Converter{}
}

// Macros returns the macro functions for Trino queries.
func (d *Datasource) Macros() sqlds.Macros {
	return sqlds.Macros{}
}

// buildDSN constructs a Trino DSN string from Grafana datasource settings.
func buildDSN(settings backend.DataSourceInstanceSettings) (string, error) {
	u, err := url.Parse(settings.URL)
	if err != nil {
		return "", fmt.Errorf("invalid URL %q: %w", settings.URL, err)
	}

	// Default to "grafana" user if no basic auth user is set
	user := "grafana"
	if settings.BasicAuthEnabled && settings.BasicAuthUser != "" {
		user = settings.BasicAuthUser
	}

	// Build the Trino DSN: http[s]://user@host:port
	dsn := fmt.Sprintf("%s://%s@%s", u.Scheme, url.PathEscape(user), u.Host)

	return dsn, nil
}
