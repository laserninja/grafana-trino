# Trino Datasource for Grafana

Connect [Grafana](https://grafana.com/) to [Trino](https://trino.io/) — a distributed SQL query engine designed for fast analytic queries against data of any size.

## Features

- Write raw SQL queries with full Trino dialect support
- Output formats: Table, Time Series, Logs
- Grafana alerting and annotations support
- Template variable support with proper SQL escaping

## Getting Started

1. Install the plugin from the Grafana Plugin Catalog
2. Add a new Trino data source in **Connections → Data sources**
3. Set the URL to your Trino coordinator (e.g., `http://trino:8080`)
4. Click **Save & test**

## Query Editor

Write SQL queries directly:

```sql
SELECT date_trunc('hour', created_at) AS time, count(*) AS value
FROM events
GROUP BY 1
ORDER BY 1
```

Choose the output format (Table, Time Series, or Logs) from the format dropdown.

## Links

- [GitHub Repository](https://github.com/laserninja/grafana-trino)
- [Issue Tracker](https://github.com/laserninja/grafana-trino/issues)
