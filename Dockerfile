# syntax=docker/dockerfile:1

## Build
FROM golang:1.25.5-bookworm AS build

WORKDIR /app

# Download Go modules
# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY petstore ./petstore
COPY main.go .

# Build statically linked binary with size optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /petstore-server

## Deploy
# Use static-debian12 for a smaller, secure base image
FROM gcr.io/distroless/static-debian12

# OpenContainers specific labels
LABEL org.opencontainers.image.title="Go gRPC Petstore Server"
LABEL org.opencontainers.image.description="A sample petstore server written in Go with gRPC and MongoDB"
LABEL org.opencontainers.image.source="https://github.com/middleland/go-servergrpc-petstore"
LABEL org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /

COPY --from=build /petstore-server /petstore-server

EXPOSE 8090

USER nonroot:nonroot

ENTRYPOINT ["/petstore-server"]
