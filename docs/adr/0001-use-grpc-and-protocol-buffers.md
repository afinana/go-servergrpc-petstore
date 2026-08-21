# ADR-0001: Use gRPC and Protocol Buffers for Petstore Service

## Status
Accepted

## Date
2026-01-20

## Context
The Petstore reference application was originally implemented as an HTTP/REST API utilizing JSON payloads over HTTP/1.1 (`go-server-petstore`). For microservice environments, high throughput, strongly-typed contracts, and binary serialization are preferred to minimize serialization overhead, enforce API schema evolution rules, and provide seamless multi-language client stub generation.

## Decision
1. Define the service contract using Protocol Buffers 3 (`proto3`) in `protobuffer/swaggerpetstore/petstore.proto`.
2. Implement a gRPC server in Go using `google.golang.org/grpc` and `google.golang.org/protobuf`.
3. Provide service implementations for `Pet`, `Store`, and `User` domains compliant with the OpenAPI Petstore 2.0 specification.

## Consequences

### Positive
- **Strong Typing & Contract Safety**: Client and server code generation eliminates contract mismatch issues.
- **Performance**: Protocol Buffers binary serialization provides smaller payloads and faster processing compared to JSON over HTTP/1.1.
- **Streaming & Concurrency**: Enables native HTTP/2 capabilities such as multiplexing and bidirectional streaming if needed in the future.

### Negative / Trade-offs
- Standard web browsers cannot call gRPC directly without gRPC-Web proxying or HTTP-to-gRPC gateways (such as grpc-gateway).
- Local debugging requires gRPC-aware CLI tools (such as `grpcurl` or Postman) rather than simple `curl`.
