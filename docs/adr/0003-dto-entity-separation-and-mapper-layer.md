# ADR-0003: Separate Transport DTOs from MongoDB Entities via Mappers

## Status
Accepted

## Date
2026-01-27

## Context
Protobuf generated structs (e.g. `petstore.Pet`, `petstore.Order`, `petstore.User`) are optimized for gRPC network transport and protobuf binary serialization. Directly using protobuf structs as MongoDB persistence models causes significant architectural friction:
- Protobuf structs contain internal fields (e.g., `state`, `sizeCache`, `unknownFields`) that bloat and corrupt MongoDB BSON documents.
- Protobuf does not define MongoDB `_id` (`primitive.ObjectID`) fields.
- Status values and dates have different representations across transport (gRPC enums, RFC3339 strings) and storage (strings, `time.Time`).

## Decision
1. Maintain separate Go structs for database entities:
   - `PetEntity`, `CategoryEntity`, `TagEntity` in `petstore/model_pet.go`, `petstore/model_category.go`, `petstore/model_tag.go`.
   - `OrderEntity` in `petstore/model_order.go`.
   - `UserEntity` in `petstore/model_user.go`.
2. Implement explicit mapper modules to convert between transport DTOs and database entities:
   - `petstore/pet_mapper.go`: `createPetEntity`, `createPetDTO`, `CreatePetListDTO`.
   - `petstore/order_mapper.go`: `createOrderEntity`, `createOrderDTO`.
   - `petstore/user_mapper.go`: `createUserEntity`, `createUserDTO`, `CreateUserListDTO`.

## Consequences

### Positive
- **Clean Architecture & Separation of Concerns**: Transport contracts and database schemas can evolve independently without leaking internal serialization concerns.
- **Data Integrity**: MongoDB documents only store clean domain fields with proper BSON tags and types (`primitive.ObjectID`, `time.Time`).
- **Maintainability**: Centralized conversion logic prevents ad-hoc field mapping across individual handler endpoints.

### Negative / Trade-offs
- Requires boilerplate mapper code and memory allocations for entity-DTO transformations during request/response handling.
