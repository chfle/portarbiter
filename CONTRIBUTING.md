# Contributing to portarbiter

portarbiter is a system administration tool.
Changes must be correct, safe, and reproducible.

## Requirements

- Linux (Debian/Ubuntu recommended)
- Go 1.22+
- Docker (for Docker-related development)
- systemd-based system

## Project Structure

- cmd/portarbiter/        – CLI entry point (minimal only)
- internal/app/           – orchestration logic
- internal/resolve/       – ownership resolution (process, systemd, docker)
- internal/detect/        – low-level port/PID detection
- internal/version/       – build-time version information
- man/                    – manual pages
- packaging/              – Debian packaging files

Do NOT put business logic into main.go.

## Versioning

Versioning is Git-tag based.

- Create a release tag:

  git tag -a v0.1.0 -m "portarbiter v0.1.0"

- The version is injected at build time.
- Do NOT hardcode versions in Go files.

## Building

Local binary:

  make build

Debian package:

  make deb

Version information:

  make version

## Safety Rules

- No feature may kill supervisor services blindly.
- Docker compose projects must require confirmation.
- SSH user sessions must never be treated as ssh.service ownership.
- Dry-run must remain the default behavior.

## Testing Expectations

Before submitting a PR, test at least:

- SSH session with forwarded ports
- Docker container with published ports
- docker-compose project
- systemd service listening on a port

## Coding Style

- Explicit logic over cleverness
- No reflection
- No runtime shell parsing unless unavoidable
- Prefer correctness over performance

## Commits

- Small, focused commits
- Imperative commit messages:
  - "resolve: fix docker published port detection"
  - "package: update deb postinst"

## Questions

If unsure, open an issue or discussion before implementing behavior-changing logic.

