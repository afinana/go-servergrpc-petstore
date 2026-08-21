# syntax=docker/dockerfile:1

## Build
FROM golang:1.25-bookworm AS build

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

# Copy repository source code including vendored dependencies
COPY . .

# Build statically linked binary offline using vendor
RUN go build -mod=vendor -ldflags="-w -s" -o /petstore-server main.go

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
