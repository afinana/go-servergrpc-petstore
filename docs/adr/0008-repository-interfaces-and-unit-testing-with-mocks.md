# ADR-0008: Repository Interfaces and Unit Testing with Mocks

## Status
Accepted

## Date
2026-08-21

## Context
Previously, the `Application` struct held direct concrete references to MongoDB repository structs (`*PetModel`, `*StoreModel`, `*UserModel`). This tightly coupled the service handlers (`AddPet`, `PlaceOrder`, `CreateUser`, etc.) directly to MongoDB. As a result:
1. Handler logic could not be verified in unit tests without running an active MongoDB instance.
2. In CI/CD or lightweight local test runs where MongoDB is absent, tests were skipped.
3. Violations of the Dependency Inversion Principle (DIP) made it difficult to introduce alternative persistence stores (such as PostgreSQL, Memory, or DynamoDB) or test error recovery paths.

## Decision
1. Define clear repository interfaces (`PetRepository`, `StoreRepository`, `UserRepository`) in `petstore/repository.go`.
2. Update `Application` to depend upon these repository interfaces rather than concrete struct pointers.
3. Ensure MongoDB models (`PetModel`, `StoreModel`, `UserModel`) implement the repository interfaces.
4. Implement in-memory mock repositories (`MockPetRepository`, `MockStoreRepository`, `MockUserRepository`) and pure unit tests for all handler methods across Pet, Store, and User APIs.

## Consequences
### Positive
- Unit tests execute instantaneously without external network or database dependencies.
- 100% test pass rate in environments without MongoDB (with integration tests executing when MongoDB is reachable).
- Full adherence to Hexagonal / Clean Architecture and SOLID principles.
- Easy to test boundary conditions, invalid inputs, and error propagation paths.

### Negative
- Additional interface definitions and mock types to maintain when adding repository methods.
