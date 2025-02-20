# Use official Golang image as a build stage
FROM golang:1.20 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application (output binary will be named `main`)
RUN go build -o main .

# Use a minimal image to run the built binary
FROM alpine:latest  

# Set working directory for the final image
WORKDIR /root/

# Copy the compiled binary from builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]
