# # Use Ubuntu as the builder image
# FROM ubuntu:24.04 AS builder

# WORKDIR /app

# # Install necessary dependencies: Git, Go, and CA certificates
# RUN apt update && apt install -y \
#     golang \
#     git \
#     ca-certificates \
#     && rm -rf /var/lib/apt/lists/*

# # Set Go environment variables
# ENV GOPROXY=https://proxy.golang.org,direct
# ENV GOSUMDB=sum.golang.org


# # Copy Go modules files and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the application source code
# COPY . .

# # Build the Go binary inside the container
# RUN GOARCH=amd64 GOOS=linux go build -o main .

# # Use Ubuntu as the final base image
# FROM ubuntu:24.04

# WORKDIR /app

# # Install CA certificates in the final image
# RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

# # Copy the compiled binary from the builder stage
# COPY --from=builder /app/main .

# # Ensure the binary has execution permissions
# RUN chmod +x main

# # Expose the Go app's port (5000)
# EXPOSE 5000

# # Run the binary
# CMD ["./main"]

# Gunakan image yang lebih ringan (Alpine)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install Git & SSL certificates
RUN apk add --no-cache git ca-certificates

# Set environment agar build lebih cepat & lebih kecil
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://proxy.golang.org,direct

# Copy go.mod dan go.sum lebih dulu agar caching lebih optimal
COPY go.mod go.sum ./
RUN go mod download

# Copy semua kode sumber setelah dependencies terinstall
COPY . .

# Build binary Go yang lebih kecil dan cepat
RUN go build -trimpath -ldflags="-s -w" -o main .

# Gunakan image runtime yang sangat kecil (Alpine)
FROM alpine:3.19

WORKDIR /app

# Install hanya sertifikat SSL untuk HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary dari tahap build
COPY --from=builder /app/main .

# Beri izin eksekusi
RUN chmod +x main

EXPOSE 5000
CMD ["./main"]
