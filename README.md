# Go gRPC Petstore Server (MongoDB)

This is a sample microservice for the Petstore application using Go, gRPC, and MongoDB. It demonstrates a basic implementation of the Petstore API using Protocol Buffers and gRPC.

## Prerequisites

- **Go**: 1.25.5 or later
- **Docker**: For running MongoDB
- **MongoDB**: (Optional if not using Docker)

## Getting Started

### 1. Start MongoDB

The application requires a MongoDB instance. You can start one easily using Docker:

```bash
docker run -d -p 27017:27017 --name test-mongo mongo:latest
```

### 2. Run the Server

To start the gRPC server (listening on port `8090`):

```bash
go run main.go
```

You should see output similar to:
```
INFO    2026/01/26 22:00:00 Database connection established
INFO    2026/01/26 22:00:00 Starting server on localhost:8090
2026/01/26 22:00:00 server listening at [::]:8090
```

### 3. Run Tests

Integration tests are located in the `petstore` package and require the MongoDB container to be running.

```bash
go test -v ./petstore/...
```

**Note:** The tests will clear the `pets` collection in the `petstore_test` database.

## Project Structure

- **`main.go`**: Entry point for the server. Handles database connection and gRPC server initialization.
- **`petstore/`**: Contains the generated protobuf code, API implementation, and models.
    - **`petstore.pb.go`**: Generated protobuf code.
    - **`petstore_grpc.pb.go`**: Generated gRPC service code.
    - **`api_*.go`**: Implementation of the gRPC service methods.
    - **`mongo_*_model.go`**: MongoDB data access layer.
    - **`*test.go`**: Integration tests.
- **`proto/`**: Contains the `petstore.pb` (protobuf definition).

## API Implementation Status

| Service | Status | Notes |
| :--- | :--- | :--- |
| **Pet** | ✅ Partial | `AddPet`, `GetPetById`, `FindPetsByStatus`, `DeletePet` implemented. |
| **Store** | ⚠️ Pending | Returns `Unimplemented` error. |
| **User** | ⚠️ Pending | Returns `Unimplemented` error. |

## Development

### Regenerating Protobuf Code

If you modify `petstore.pb`, you can regenerate the Go code:

1.  Install `protoc` and the Go plugins:
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
2.  Run `protoc`:
    ```bash
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        petstore/petstore.pb
    ```

## License

This project is generated from the [Swagger Petstore](http://petstore.swagger.io) API definition.
