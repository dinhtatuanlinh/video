# Stage 1: Build the Go binary
FROM golang:1.22.1 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download
# Copy the migration folder
COPY db/migration /app/db/migration
# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Create a minimal container
FROM alpine:latest

# Install CA certificates for outbound HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app /app
COPY dev.env .
# Command to run the app
CMD ["./app"]
