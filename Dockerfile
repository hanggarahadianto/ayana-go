# 1️⃣ Use official Go image to build
FROM golang:1.21-alpine AS builder

# 2️⃣ Set the working directory inside the container
WORKDIR /app/ayana-go

# 3️⃣ Copy go.mod and go.sum files
COPY ayana-go/go.mod ayana-go/go.sum ./

# 4️⃣ Download dependencies
RUN go mod download

# 5️⃣ Copy the rest of the application files
COPY ayana-go .

# 6️⃣ Build the Go application and place it in ayana-go
RUN go build -o /app/ayana-go/main

# 7️⃣ Use a lightweight image for runtime
FROM alpine:latest

# 8️⃣ Set working directory in final container
WORKDIR /app/ayana-go

# 9️⃣ Copy the built binary from the builder stage
COPY --from=builder /app/ayana-go/main .

# 🔟 Ensure the binary has execution permissions
RUN chmod +x /app/ayana-go/main

# 1️⃣1️⃣ Expose the application's port
EXPOSE 8080

# 1️⃣2️⃣ Start the application
CMD ["/app/ayana-go/main"]
