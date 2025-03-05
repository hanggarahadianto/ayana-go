# Use Ubuntu as the builder image
# FROM ubuntu:24.04 AS builder
FROM golang:1.21 AS builder

WORKDIR /app

# Install necessary dependencies: Git, Go, and CA certificates
RUN apt update && apt install -y \
    golang \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set Go environment variables
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org


# Copy Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary inside the container
RUN GOARCH=amd64 GOOS=linux go build -o main .

# Use Ubuntu as the final base image
FROM ubuntu:24.04


WORKDIR /app

# Install CA certificates in the final image
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary has execution permissions
RUN chmod +x main

# Expose the Go app's port (5000)
EXPOSE 5000

# Run the binary
CMD ["./main"]