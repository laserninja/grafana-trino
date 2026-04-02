# Contributing

Thank you for considering contributing to the Grafana Trino Datasource Plugin!

## Getting Started

1. Fork the repository
2. Clone your fork and set up the dev environment (see [DEVELOPMENT.md](DEVELOPMENT.md))
3. Create a feature branch: `git checkout -b my-feature`
4. Make your changes
5. Submit a pull request

## Requirements for Pull Requests

- **Tests required**: Every feature or bug fix must include tests (unit and/or E2E)
- **CI must pass**: Lint, typecheck, build, and all tests must be green
- **Keep changes focused**: One logical change per PR
- **Update docs**: If your change affects usage, update README.md and relevant docs

## Code Style

- **Frontend**: ESLint + Prettier (run `npm run lint:fix`)
- **Backend**: Standard Go formatting (`gofmt`)
- Follow existing patterns in the codebase

## Reporting Issues

Please use [GitHub Issues](https://github.com/laserninja/grafana-trino/issues) to report bugs or request features. Include:

- Grafana version
- Trino version
- Steps to reproduce
- Expected vs actual behavior

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
