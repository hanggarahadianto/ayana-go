# Use Ubuntu as the base image for building
FROM ubuntu:24.04 AS builder

# Install necessary dependencies
RUN apt update && apt install -y golang

# Set the working directory
WORKDIR /app

# Copy Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN GOARCH=amd64 GOOS=linux go build -o main .

# Use Ubuntu as the minimal base image for the final container
FROM ubuntu:24.04

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary has execution permissions
RUN chmod +x main


EXPOSE 8080

# Run the binary
CMD ["./main"]
