# ADR-0005: Multi-Stage Docker Builds with Distroless Base Images

## Status
Accepted

## Date
2026-01-27

## Context
Production container images for microservices should have a minimal attack surface, fast transfer times, and zero extraneous dependencies (e.g. package managers, shells, build tools).

## Decision
1. **Multi-Stage Build**:
   - Build stage: `golang:1.25.5-bookworm` compiles a statically-linked binary with stripped debug info (`CGO_ENABLED=0`, `-ldflags="-w -s"`).
   - Deployment stage: `gcr.io/distroless/static-debian12`.
2. **Security Best Practices**:
   - Run as non-root user (`USER nonroot:nonroot`).
   - Add OpenContainers metadata labels (`org.opencontainers.image.*`).
   - Use `.dockerignore` to exclude local binaries, git directories, IDE files, and test files from the build context.

## Consequences

### Positive
- **Minimal Image Size**: Final container image is very small (~20-30MB).
- **Hardened Security**: Absence of shell and package managers eliminates large categories of container breakout and supply chain vulnerabilities.
- **Fast Deployments**: Reduced bandwidth and storage requirements across CI/CD and container registries.

### Negative / Trade-offs
- Executing shell commands (`docker exec -it ... /bin/sh`) inside the running container is not possible for interactive debugging; diagnostic logs must rely on application output.
