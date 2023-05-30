# Issues

1. File type is not recognized->changed to string (base64).
2. schema:string is not recognized in Login op. Changed to ApiResponse.


# Command

```console
openapi2proto -spec petstore-openapi.json > petstore.pb
```

# From protobuffer to go code

Add this option to petstore.pb file.

option go_package = "./";


```console
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative petstore.pb
```