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

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /petstore-server

## Deploy
# Use static-debian12 for a smaller, secure base image
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=build /petstore-server /petstore-server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/petstore-server"]
