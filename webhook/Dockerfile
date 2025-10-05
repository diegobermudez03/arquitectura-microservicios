# Start from the official Golang image
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o cards-app main.go

# Use a minimal image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/cards-app .

# Expose the application port (change 8080 if your app uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./cards-app"]