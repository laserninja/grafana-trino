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
- E2E tests with Playwright and `@grafana/plugin-e2e`
- CI/CD with GitHub Actions (build, lint, test, E2E, release)
- Docker Compose development environment with Trino
