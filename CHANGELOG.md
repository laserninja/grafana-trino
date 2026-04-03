# Changelog

## 2.0.0 (Unreleased)

Complete rewrite using the latest Grafana plugin SDK and tooling.

### Breaking Changes

- Minimum Grafana version is now 10.0.0
- Plugin rebuilt from scratch — configuration is compatible but internal implementation has changed

### Added

- Modern backend using `grafana-plugin-sdk-go` v0.291.0 and `sqlds/v4`
- Configuration UI built with `@grafana/plugin-ui` components
- SQL query editor with syntax highlighting
- Output format selection: Table, Time Series, Logs
- All 7 authentication methods: Basic auth, JWT access token, OAuth2 client credentials, TLS/mTLS, user impersonation, roles, client tags
- All 8 query macros: `$__timeFrom`, `$__timeTo`, `$__timeGroup`, `$__timeFilter`, `$__dateFilter`, `$__unixEpochFilter`, `$__unixEpochGroup`, `$__parseTime`
- 15+ type converters for Trino SQL types (varchar, integer, timestamp, boolean, etc.)
- Template variable support with proper SQL escaping (single-value, multi-value)
- Per-instance OAuth2 token cache (no global singleton)
- Standard variable support for dashboard variable queries
- E2E tests with Playwright and `@grafana/plugin-e2e`
- CI/CD with GitHub Actions (build, lint, test, E2E, release)
- Docker Compose development environment with Trino
- Comprehensive backend unit tests (42+ tests)
