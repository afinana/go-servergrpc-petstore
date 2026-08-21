# ADR-0007: gRPC Server Reflection and Standard Health Checks

## Status
Accepted

## Date
2026-08-21

## Context
When running gRPC services in modern cloud and container ecosystems (such as Kubernetes or Docker Compose):
1. **Dynamic discovery & CLI inspection**: Developers and operators interacting with the gRPC service via CLI tools like `grpcurl` or GUI tools like Postman previously had to supply the `.proto` definition files locally.
2. **Liveness & Readiness Probing**: Container orchestrators (Kubernetes, Nomad, Docker Swarm) require standardized health check endpoints (`grpc.health.v1.Health`) to detect degraded instances, manage rolling rollouts, and ensure traffic is only routed to healthy pods.

## Decision
1. Register standard gRPC Reflection (`google.golang.org/grpc/reflection`) on the server during startup.
2. Implement and register the official gRPC Health Checking Protocol (`google.golang.org/grpc/health/grpc_health_v1`).
3. Set serving statuses to `SERVING` for both the default service `""` and `"petstore.SwaggerPetstoreService"`.
4. Transition health check statuses to `NOT_SERVING` upon receipt of termination signals (`SIGINT`/`SIGTERM`) prior to executing `server.GracefulStop()`.

## Consequences
### Positive
- `grpcurl` and API testing clients can discover service methods, requests, and responses without manually loading `.proto` files (`grpcurl -plaintext localhost:8090 list`).
- Kubernetes / container health probes can natively monitor container status via `grpc-health-probe` or native gRPC probes.
- Zero-downtime rolling updates are enabled by failing health probes before in-flight requests finish draining.

### Negative
- Minimal binary size overhead for reflection schema registration.
