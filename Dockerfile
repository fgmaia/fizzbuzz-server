# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fizzbuzz-server ./cmd/fizzbuzz-server

# Stage 2: Create a minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/fizzbuzz-server .
# Copy any config files if needed
COPY --from=builder /app/.env* ./

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./fizzbuzz-server"]