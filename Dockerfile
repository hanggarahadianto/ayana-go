


# Use the official Golang image
FROM golang:1.20-alpine AS builder

# Install necessary dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal Alpine Linux image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary is executable
RUN chmod +x /app/main

# Expose the port that your application listens on
EXPOSE 8080

# Run the application
CMD ["/app/main"]

