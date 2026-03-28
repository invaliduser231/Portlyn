# Contributing

Thank you for contributing to Portlyn.

## Before you start

- Open an issue for significant changes before starting implementation.
- Keep changes focused. Avoid mixing refactors, new features, and cleanup in one pull request.
- If the change affects behavior, include tests or explain why tests are not needed.
- If the change affects setup, configuration, or APIs, update the documentation in the same pull request.

## Development expectations

- Backend changes should pass:
  - `go test ./...`
  - `go build ./...`
- Configuration changes should be validated with:
  - `go run ./cmd/configcheck`
- Docker or deployment changes should keep `docker-compose.yml` and `.env.docker.example` aligned.

## Pull requests

- Use a clear title and describe the user-visible impact.
- Include screenshots for UI changes when relevant.
- Call out breaking changes explicitly.
- Keep PRs reviewable. Smaller PRs merge faster and with less risk.

## Scope

We may decline contributions that:

- expand the product far beyond its current direction
- introduce large dependency or maintenance costs
- weaken security defaults
- conflict with the commercial licensing direction of the project

