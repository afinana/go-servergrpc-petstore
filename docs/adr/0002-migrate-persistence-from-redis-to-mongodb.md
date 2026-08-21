# ADR-0002: Migrate Persistence from Redis to MongoDB

## Status
Accepted

## Date
2026-01-27

## Context
Early prototypes of the service utilized Redis as a key-value and JSON cache store (`redis_pet_model.go`, `redis_order_model.go`, `redis_user_model.go`). However, the Petstore domain model requires:
- Complex document schemas with nested structures (e.g. `Category`, `Tags`, `PhotoUrls`).
- Multi-field querying capabilities (e.g. filtering pets by status `$or` arrays, searching pets by embedded tags using `$elemMatch`).
- Aggregation pipelines (e.g. calculating pet inventory grouped by status).

## Decision
1. Replace Redis models with MongoDB document repositories (`mongo_pet_model.go`, `mongo_store_model.go`, `user_model.go`).
2. Utilize the official Go MongoDB Driver (`go.mongodb.org/mongo-driver`).
3. Store documents in three distinct collections: `pets`, `stores`, and `users`.
4. Implement MongoDB aggregation pipelines for inventory counts and query filters (`$or`, `$elemMatch`).

## Consequences

### Positive
- **Rich Querying**: Native support for filtering by nested arrays, status lists, and username queries without manual indexing layers.
- **Aggregation Framework**: Simplified inventory status calculations via `$group` and `$sum` aggregation pipelines.
- **Document Flexibility**: BSON mapping natively accommodates dynamic nested objects (`CategoryEntity`, `TagEntity`).

### Negative / Trade-offs
- Requires a MongoDB database instance in development and test environments instead of lightweight in-memory Redis.
- Higher baseline memory footprint compared to simple key-value lookups.
