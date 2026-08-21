# gRPC API Guide & Reference

This guide details all available gRPC procedures exposed by the `SwaggerPetstoreService` in `protobuffer/swaggerpetstore/petstore.proto`.

## Service Overview

Package: `petstore`  
Service: `SwaggerPetstoreService`  
Default Port: `8090`

---

## Pet API

### 1. `AddPet`
Adds a new pet to the database.

- **RPC**: `rpc AddPet(AddPetRequest) returns (google.protobuf.Empty)`
- **Request**:
  ```protobuf
  message AddPetRequest {
    Pet body = 1;
  }
  ```
- **Error Codes**: `InvalidArgument`, `Internal`

### 2. `GetPetById`
Fetches a single pet by its numerical ID.

- **RPC**: `rpc GetPetById(GetPetByIdRequest) returns (Pet)`
- **Request**:
  ```protobuf
  message GetPetByIdRequest {
    int64 petId = 1;
  }
  ```
- **Error Codes**: `NotFound`, `Internal`

### 3. `FindPetsByStatus`
Finds pets matching one or more status enum values (`STATUS_AVAILABLE`, `STATUS_PENDING`, `STATUS_SOLD`).

- **RPC**: `rpc FindPetsByStatus(FindPetsByStatusRequest) returns (FindPetsByStatusResponse)`
- **Request**:
  ```protobuf
  message FindPetsByStatusRequest {
    repeated Status status = 1;
  }
  ```

### 4. `FindPetsByTags`
Finds pets matching one or more tag names.

- **RPC**: `rpc FindPetsByTags(FindPetsByTagsRequest) returns (FindPetsByTagsResponse)`
- **Request**:
  ```protobuf
  message FindPetsByTagsRequest {
    repeated string tags = 1;
  }
  ```

### 5. `UpdatePet`
Updates an existing pet in the store.

- **RPC**: `rpc UpdatePet(UpdatePetRequest) returns (google.protobuf.Empty)`

### 6. `UpdatePetWithForm`
Updates a pet's name and/or status using form parameters.

- **RPC**: `rpc UpdatePetWithForm(UpdatePetWithFormRequest) returns (google.protobuf.Empty)`

### 7. `DeletePet`
Deletes a pet by its numerical ID.

- **RPC**: `rpc DeletePet(DeletePetRequest) returns (google.protobuf.Empty)`

### 8. `UploadFile`
Uploads an image/file metadata for a pet.

- **RPC**: `rpc UploadFile(UploadFileRequest) returns (ApiResponse)`

---

## Store API

### 1. `PlaceOrder`
Places an order for a pet.

- **RPC**: `rpc PlaceOrder(PlaceOrderRequest) returns (Order)`
- **Request**:
  ```protobuf
  message PlaceOrderRequest {
    Order body = 1;
  }
  ```

### 2. `GetOrderById`
Fetches an order by its numerical ID.

- **RPC**: `rpc GetOrderById(GetOrderByIdRequest) returns (Order)`
- **Error Codes**: `NotFound`, `Internal`

### 3. `DeleteOrder`
Deletes an order by its numerical ID.

- **RPC**: `rpc DeleteOrder(DeleteOrderRequest) returns (google.protobuf.Empty)`

### 4. `GetInventory`
Returns a breakdown of pet counts grouped by status.

- **RPC**: `rpc GetInventory(google.protobuf.Empty) returns (GetInventoryResponse)`
- **Response**:
  ```protobuf
  message GetInventoryResponse {
    repeated MapInventory items = 1;
  }
  message MapInventory {
    string name = 1;
    int64 value = 2;
  }
  ```

---

## User API

### 1. `CreateUser`
Creates a new user record.

- **RPC**: `rpc CreateUser(CreateUserRequest) returns (google.protobuf.Empty)`

### 2. `CreateUsersWithArrayInput` / `CreateUsersWithListInput`
Batch-creates multiple user records from a list.

- **RPC**: `rpc CreateUsersWithArrayInput(CreateUsersWithArrayInputRequest) returns (google.protobuf.Empty)`
- **RPC**: `rpc CreateUsersWithListInput(CreateUsersWithListInputRequest) returns (google.protobuf.Empty)`

### 3. `GetUserByName`
Fetches a user by their unique username.

- **RPC**: `rpc GetUserByName(GetUserByNameRequest) returns (User)`
- **Error Codes**: `NotFound`, `Internal`

### 4. `UpdateUser`
Updates an existing user's details.

- **RPC**: `rpc UpdateUser(UpdateUserRequest) returns (google.protobuf.Empty)`

### 5. `DeleteUser`
Deletes a user by username.

- **RPC**: `rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty)`

### 6. `LoginUser`
Authenticates a user by username and password.

- **RPC**: `rpc LoginUser(LoginUserRequest) returns (ApiResponse)`
- **Error Codes**: `Unauthenticated`, `Internal`

### 7. `LogoutUser`
Terminates a user session.

- **RPC**: `rpc LogoutUser(google.protobuf.Empty) returns (google.protobuf.Empty)`

---

## Interacting via `grpcurl`

With **gRPC Server Reflection** enabled, you do not need to provide proto files:

```bash
# List all registered services (via reflection)
grpcurl -plaintext localhost:8090 list

# Describe a specific service schema
grpcurl -plaintext localhost:8090 describe petstore.SwaggerPetstoreService

# Add a pet
grpcurl -plaintext -d '{"body": {"id": 1001, "name": "Doggie", "status": "STATUS_AVAILABLE"}}' localhost:8090 petstore.SwaggerPetstoreService/AddPet

# Get pet by ID
grpcurl -plaintext -d '{"petId": 1001}' localhost:8090 petstore.SwaggerPetstoreService/GetPetById

# Get inventory
grpcurl -plaintext localhost:8090 petstore.SwaggerPetstoreService/GetInventory

# Check service health (Standard gRPC Health Protocol)
grpcurl -plaintext -d '{"service": "petstore.SwaggerPetstoreService"}' localhost:8090 grpc.health.v1.Health/Check
```
