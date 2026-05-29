# Security Policy

## Supported versions

Security fixes are applied on a best-effort basis to the latest maintained code in the default branch.

## Reporting a vulnerability

Please do not open public GitHub issues for suspected vulnerabilities.

Report security issues privately through either:

- Email: `security@portlyn.dev`
- GitHub Private Vulnerability Reporting: use the **Security** tab of this repository

Include the following:

- a clear description of the issue
- affected version or commit
- reproduction steps or proof of concept
- impact assessment if known

## Disclosure policy

- We will acknowledge valid reports as quickly as practical.
- We may ask for reproduction details or environment information.
- Please allow time for investigation and remediation before public disclosure.

## Telemetry

Portlyn collects no telemetry.
There is no analytics SDK, no phone-home, and no automatic update check.
The only outbound traffic comes from features you explicitly configure
(ACME, webhooks, OIDC, DNS provider APIs, CrowdSec LAPI).
