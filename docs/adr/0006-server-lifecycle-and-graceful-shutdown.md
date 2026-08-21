# ADR-0006: Unified Configuration and Graceful Server Shutdown

## Status
Accepted

## Date
2026-02-20

## Context
Production servers running in orchestrators (e.g. Kubernetes, Docker Compose, Nomad) send termination signals (`SIGINT`, `SIGTERM`) during rolling updates and scaling operations. If the server terminates abruptly:
- In-flight gRPC RPC calls are terminated unexpectedly.
- Open database connections and connection pools are dropped ungracefully.

Additionally, configuration flags needed to be simplified and aligned with the REST implementation (`go-server-petstore`), removing redundant flags.

## Decision
1. **Unified Network Address Flag**: Consolidate `serverAddr` and `serverPort` into a single `-serverAddr` flag (e.g. `localhost:8090`, `:8090`, `0.0.0.0:8090`).
2. **Graceful Shutdown**: Intercept `os.Interrupt` and `syscall.SIGTERM` signals, executing `server.GracefulStop()` to finish active RPCs before exiting.
3. **Dedicated Disconnect Context**: Use a dedicated timeout context (`context.WithTimeout(context.Background(), 10*time.Second)`) for `client.Disconnect()` rather than relying on the expired startup context.

## Consequences

### Positive
- **Zero-Downtime Deployments**: In-flight requests complete normally before the process terminates.
- **Clean Resource Release**: MongoDB connection pools and network listeners are closed orderly.
- **Consistent Configuration**: Matches the single address specification pattern used in modern Go services.

### Negative / Trade-offs
- Shutdown sequence may delay container termination up to the orchestrator grace period if an in-flight RPC is long-running.
