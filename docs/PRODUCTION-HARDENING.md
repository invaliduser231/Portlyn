# Production Hardening Guide

This project ships with a developer-friendly default Docker Compose stack. That default is not a
production security posture by itself.

## Required before internet exposure

- Set `ALLOW_INSECURE_DEV_MODE=false`.
- Set strong values for `JWT_SECRET`, `POSTGRES_PASSWORD`, and `GRAFANA_ADMIN_PASSWORD`.
- Disable bootstrap login after the first admin setup: `BOOTSTRAP_ADMIN_ENABLED=false`.
- Enroll TOTP for every admin, then set `REQUIRE_MFA_FOR_ADMINS=true`.
- Set `FRONTEND_BASE_URL`, `CORS_ALLOWED_ORIGINS`, and OIDC redirect settings to the final public URL.
- Enable TLS for the shared admin/proxy entrypoint before exposing it publicly.

## Exposure model

- Treat the admin UI and API as privileged surfaces.
- Do not expose Grafana, Loki, or `/metrics` publicly without an explicit access layer.
- Prefer a private network, VPN, or upstream access control in front of observability endpoints.
- If you publish the main proxy entrypoint, separate public service routes from operator-only endpoints.

## Authentication and secrets

- Keep `OTP_RESPONSE_INCLUDES_CODE=false` outside local development.
- Use a dedicated SMTP account for OTP and route email flows.
- Rotate credentials when operators change.
- Back up secrets separately from the application host.

## Storage and recovery

- Back up PostgreSQL on a schedule and test restore regularly.
- Back up `/data` when using local certificate storage.
- Record the exact application version used for each backup.

## Release and upgrade discipline

- Run `go test ./...`, `go vet ./...`, and the frontend build before every release.
- Review the security workflow results before tagging a release.
- Keep a rollback path for database and config changes.
