# Use the official Golang base image with Alpine Linux
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application source code into the container
COPY . .

# Build the Go application
RUN go build -o main

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

# Set the environment variable to use .env file
ENV GIN_MODE=release

# Command to run the application
CMD ["/app/main"]
