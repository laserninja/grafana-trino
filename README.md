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
4. Configure authentication if needed (Basic Auth, TLS)
5. Click **Save & test**

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

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for local development setup.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
