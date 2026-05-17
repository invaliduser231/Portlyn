# Release Process

## Scope

This repository does not yet have an automated publish pipeline. Until that exists, use the
following release checklist for tagged versions.

## Pre-release checklist

- Confirm CI is green on the release commit.
- Review the latest security workflow results.
- Run `go run ./cmd/configcheck`.
- Run `go test ./...`.
- Run the frontend build from `frontend/`.
- Confirm `docker compose --env-file .env.docker.example config` still succeeds.
- Update `CHANGELOG.md`.
- Add upgrade notes for any config, auth, routing, or storage changes.

## Tagging

- Use semantic versions: `vMAJOR.MINOR.PATCH`.
- Tag only commits already merged to `main`.
- Keep the tag message short and reference the changelog section.

## Release notes

Each release should call out:

- operator-visible changes
- security-impacting changes
- config or environment variable changes
- migration or rollback notes
- known limitations that still apply

## Post-release checks

- Verify the tag points at the intended commit.
- Verify the Dockerfiles still build from the tagged source.
- Open follow-up issues for deferred fixes instead of carrying undocumented release debt.
