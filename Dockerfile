


# Use an official Golang image to build the app
FROM golang:1.21 AS builder

WORKDIR /app

# Copy Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary inside the container
RUN go build -o main .

# Use a minimal base image for the final container
FROM alpine:latest  

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary has execution permissions
RUN chmod +x main  

# Run the binary
CMD ["./main"]
