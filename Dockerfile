# Multi-stage Docker build for efficient images

# Build stage
FROM golang:1.25 AS build
WORKDIR /src

# Copy go mod files first (for better Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /edgelink ./cmd/edgelink

# Runtime stage - minimal image for security
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# Copy binary and config from build stage
COPY --from=build /edgelink /usr/local/bin/edgelink
COPY configs/example.yaml /app/config.yaml

# Expose HTTP port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/usr/local/bin/edgelink", "--config", "/app/config.yaml"]
