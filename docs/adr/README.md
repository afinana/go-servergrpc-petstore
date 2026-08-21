# Architecture Decision Records (ADRs)

This directory documents the key architectural decisions made during the evolution of the **Go gRPC Petstore Server** project.

Each record captures the context, decision, and consequences of a significant architectural choice.

## Statuses
- **Accepted**: Decision made and implemented in codebase.
- **Superseded**: Replaced by a subsequent decision.
- **Deprecated**: No longer relevant.

## Index of ADRs

| ADR | Title | Status | Date |
| :--- | :--- | :--- | :--- |
| [ADR-0001](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0001-use-grpc-and-protocol-buffers.md) | Use gRPC and Protocol Buffers for Petstore Service | Accepted | 2026-01-20 |
| [ADR-0002](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0002-migrate-persistence-from-redis-to-mongodb.md) | Migrate Persistence from Redis to MongoDB | Accepted | 2026-01-27 |
| [ADR-0003](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0003-dto-entity-separation-and-mapper-layer.md) | Separate Transport DTOs from MongoDB Entities via Mappers | Accepted | 2026-01-27 |
| [ADR-0004](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0004-context-propagation-and-structured-logging.md) | Standardize Context Propagation and Colored Logging | Accepted | 2026-01-31 |
| [ADR-0005](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0005-distroless-container-and-multi-stage-builds.md) | Multi-Stage Docker Builds with Distroless Base Images | Accepted | 2026-01-27 |
| [ADR-0006](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0006-server-lifecycle-and-graceful-shutdown.md) | Unified Configuration and Graceful Server Shutdown | Accepted | 2026-02-20 |
| [ADR-0007](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0007-grpc-reflection-and-health-checks.md) | gRPC Server Reflection and Standard Health Checks | Accepted | 2026-08-21 |
| [ADR-0008](file:///home/afinana/middleland/go-servergrpc-petstore/docs/adr/0008-repository-interfaces-and-unit-testing-with-mocks.md) | Repository Interfaces and Unit Testing with Mocks | Accepted | 2026-08-21 |
