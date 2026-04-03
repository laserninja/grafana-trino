# Grafana Trino Datasource Plugin

A [Grafana](https://grafana.com/) datasource plugin for [Trino](https://trino.io/) — a distributed SQL query engine designed for fast analytic queries against data of any size.

## Features

- Connect Grafana to any Trino cluster
- Write raw SQL queries with full Trino dialect support
- Output formats: Table, Time Series, Logs
- Grafana alerting support
- Annotations support
- Template variable support with proper escaping

## Requirements

- Grafana ≥ 10.0.0
- Trino cluster accessible from the Grafana server

## Installation

### From Grafana Plugin Catalog

Search for **Trino** in Grafana → Administration → Plugins, then click **Install**.

### Manual

Download the latest release from the [GitHub releases page](https://github.com/laserninja/grafana-trino/releases), extract into your Grafana plugins directory, and restart Grafana.

## Configuration

1. Go to **Connections → Data sources → Add data source**
2. Search for **Trino**
3. Set the **URL** to your Trino coordinator (e.g., `http://trino:8080`)
4. Configure authentication if needed (see below)
5. Click **Save & test**

### Authentication

| Method | Description |
|--------|-------------|
| **Basic Auth** | Username/password via Grafana's built-in basic auth settings |
| **Access Token** | Static JWT token for direct Trino authentication |
| **OAuth2 Client Credentials** | Automatic token retrieval via client credentials flow (Token URL, Client ID, Client Secret) |
| **TLS / mTLS** | CA certificate, client certificate, and client key via Grafana's TLS settings |
| **User Impersonation** | Sets `X-Trino-User` to the current Grafana user's login |
| **Roles** | Catalog-specific authorization roles (e.g., `system:admin;catalog1:reader`) |
| **Client Tags** | Comma-separated tags for Trino resource group identification |

### OAuth2 Configuration

For OAuth2 client credentials flow, configure under **OAuth2 Trino Authentication**:
- **Token URL**: Your identity provider's token endpoint
- **Client ID**: OAuth2 client identifier
- **Client secret**: OAuth2 client secret (stored securely)
- **Impersonation user**: Optional user to impersonate in Trino via OAuth

## Usage

### Query Editor

Write SQL queries directly in the code editor:

```sql
SELECT date_trunc('hour', created_at) AS time, count(*) AS value
FROM events
WHERE created_at > TIMESTAMP '2024-01-01'
GROUP BY 1
ORDER BY 1
```

### Output Formats

- **Table**: Returns results as-is in table format
- **Time Series**: Expects a `time` column and one or more value columns
- **Logs**: Expects a `time` column and a message column

### Query Macros

| Macro | Description | Example Output |
|-------|-------------|----------------|
| `$__timeFrom()` | Lower time boundary | `TIMESTAMP '2024-01-01 00:00:00'` |
| `$__timeTo()` | Upper time boundary | `TIMESTAMP '2024-01-02 00:00:00'` |
| `$__timeFilter(col)` | Time range filter | `col BETWEEN TIMESTAMP '...' AND TIMESTAMP '...'` |
| `$__dateFilter(col)` | Date range filter | `col BETWEEN date '...' AND date '...'` |
| `$__timeGroup(col, interval)` | Time bucketing | `FROM_UNIXTIME(FLOOR(TO_UNIXTIME(col)/3600)*3600)` |
| `$__unixEpochFilter(col)` | Unix epoch range filter | `col BETWEEN 1704067200 AND 1704153600` |
| `$__unixEpochGroup(col, interval)` | Unix epoch bucketing | `FROM_UNIXTIME(FLOOR(col/300)*300)` |
| `$__parseTime(col, format)` | Parse time with format | `parse_datetime(col,'yyyy-MM-dd')` |

### Template Variables

Template variables are supported in queries. Values are automatically escaped:
- **Single-value**: Single quotes are escaped (`value'` → `value''`)
- **Multi-value**: Comma-separated quoted list (`'val1','val2','val3'`)

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for local development setup.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
