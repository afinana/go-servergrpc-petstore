# Go API Server GRPC- REDIS db version for swagger (1.0.1)

This is a sample microservice of Petstore application.  You can find out more about Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).  For this sample, you can use the api key `special-key` to test the authorization filters.

## Overview
This server was generated partially by the [swagger-codegen]
(https://github.com/swagger-api/swagger-codegen) project.  
By using the [OpenAPI-Spec](https://github.com/OAI/OpenAPI-Specification) from a remote server, you can easily generate a server stub.  
-

To see how to make this your own, look here:

[README](https://github.com/swagger-api/swagger-codegen/blob/master/README.md)

- API version: 1.0.6
- Build date: 2023-06-15T23:27:07.696Z


## Running the server
To run the server, follow these simple steps:

```
go run main.go
```
## Running with Docker

``` sh
docker -t go-petstore build .
docker run --name go-petstore  -p 8090:8080  go-server-petstore
```

##  How to convert Openapi specification (json file format) to protobuf by `openapi2proto`.

### Download openapi2proto and execute the following command:

```console
openapi2proto -spec swagger.json > petstore.pb
```

### From protobuffer to go code

1. Add this option to petstore.pb file.

```
option go_package = "./";
```

2. Execute the following command:

```console
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative petstore.pb
```

3. Convertion Issues:

- File type was not recognized and has been changed to string (base64).
- schema:string is not recognized on Login operation and has been changed to ApiResponse.


## Running REDIS with Docker

``` sh
docker run --name redis-stack -d -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```
