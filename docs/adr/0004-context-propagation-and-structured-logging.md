# ADR-0004: Standardize Context Propagation and Colored Logging

## Status
Accepted

## Date
2026-01-31

## Context
In microservice architectures, request lifecycles must be controllable via context deadlines, timeouts, and cancellations. Additionally, local developer experience and log parsing benefit from structured, colorized terminal logs that distinguish severity levels without adding heavy external logging frameworks.

## Decision
1. **Context Propagation**: Pass incoming gRPC `ctx context.Context` through all service handler layers down to the MongoDB database client operations (`Find`, `FindOne`, `InsertOne`, `UpdateOne`, `DeleteOne`, `Aggregate`).
2. **Colored Logger**: Implement a zero-dependency ANSI color writer (`ColoredWriter`) in `petstore/logger.go` with distinct log level formatting (`INFO` in Green, `ERROR` in Red, `DEBUG` in Blue).

## Consequences

### Positive
- **Resource Management**: When a gRPC client disconnects or times out, the active MongoDB query is cancelled immediately.
- **Observability**: Clear, color-coded terminal output during local development and debugging without third-party logging overhead.
- **Traceability**: Stack traces on errors are written via `serverError(err)` using `runtime/debug.Stack()`.

### Negative / Trade-offs
- Standard ANSI escape codes must be stripped or handled if ingested by log aggregators that do not support color formatting.
