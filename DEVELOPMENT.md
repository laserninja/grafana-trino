# Development

## Prerequisites

- [Node.js](https://nodejs.org/) >= 22
- [Go](https://golang.org/) >= 1.24
- [Mage](https://magefile.org/)
- [Docker](https://www.docker.com/) and Docker Compose

## Quick Start

```bash
# Install frontend dependencies
npm install

# Build frontend in watch mode
npm run dev

# Build backend (all platforms)
mage -v buildAll

# Start Grafana + Trino with Docker Compose
npm run server
```

Open [http://localhost:3000](http://localhost:3000) (default credentials: admin/admin). The Trino datasource is auto-provisioned.

## Project Structure

```
pkg/
  main.go                  # Backend plugin entry point
  trino/
    datasource.go          # sqlds.Driver implementation (Connect, Converters, Macros)
    datasource_test.go     # Go unit tests
src/
  module.ts                # Frontend plugin entry point
  datasource.ts            # DataSourceWithBackend implementation
  types.ts                 # TypeScript types (TrinoQuery, TrinoDataSourceOptions)
  components/
    ConfigEditor.tsx       # Datasource configuration UI
    QueryEditor.tsx        # SQL query editor with format selector
tests/
  configEditor.spec.ts     # E2E tests for configuration
  queryEditor.spec.ts      # E2E tests for query editor
provisioning/
  datasources/
    datasources.yml        # Auto-provisioned Trino datasource for dev
```

## Commands

| Command | Description |
|---------|-------------|
| `npm install` | Install frontend dependencies |
| `npm run dev` | Build frontend in watch mode |
| `npm run build` | Production frontend build |
| `npm run test:ci` | Run Jest unit tests |
| `npm run lint` | Run ESLint |
| `npm run typecheck` | Run TypeScript type checking |
| `npm run e2e` | Run Playwright E2E tests |
| `npm run server` | Start Docker Compose (Grafana + Trino) |
| `mage -v buildAll` | Build backend for all platforms |
| `go test ./pkg/...` | Run Go unit tests |

## Running E2E Tests

E2E tests require Grafana and Trino running via Docker Compose:

```bash
# Start the environment
npm run server

# In another terminal, run E2E tests
npm run e2e
```

To test against a specific Grafana version:

```bash
GRAFANA_VERSION=10.0.0 npm run server
```

## Backend Development

The backend uses [sqlds](https://github.com/grafana/sqlds) — Grafana's shared SQL datasource framework. The plugin implements the `sqlds.Driver` interface:

- `Connect()` — Opens a connection to Trino
- `Converters()` — Maps Trino SQL types to Grafana data types
- `Macros()` — Provides Trino-specific query macros
- `Settings()` — Returns driver configuration (fill mode, timeouts)
