package trino

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/sqlds/v4"
)

// Datasource implements the sqlds.Driver interface for Trino.
type Datasource struct{}

// Compile-time interface checks.
var (
	_ sqlds.Driver         = (*Datasource)(nil)
	_ sqlds.QueryArgSetter = (*Datasource)(nil)
)

// Settings returns the driver settings for Trino.
func (d *Datasource) Settings(_ context.Context, _ backend.DataSourceInstanceSettings) sqlds.DriverSettings {
	return sqlds.DriverSettings{
		FillMode: &fillModeNull,
	}
}

var fillModeNull = data.FillMissing{Mode: data.FillModeNull}

// Connect opens a SQL connection to Trino using the datasource settings.
func (d *Datasource) Connect(_ context.Context, config backend.DataSourceInstanceSettings, _ json.RawMessage) (*sql.DB, error) {
	settings, err := loadSettings(config)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	db, err := openDB(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to open Trino connection: %w", err)
	}

	return db, nil
}

// Converters returns the type converters for Trino SQL types.
func (d *Datasource) Converters() []sqlutil.Converter {
	return converters()
}

// Macros returns the macro functions for Trino queries.
func (d *Datasource) Macros() sqlds.Macros {
	return macros
}

// SetQueryArgs injects Trino-specific query arguments (user impersonation, access token, client tags)
// into the SQL query context based on the current request headers and context values.
func (d *Datasource) SetQueryArgs(ctx context.Context, headers http.Header) []interface{} {
	return setQueryArgs(ctx)
}
