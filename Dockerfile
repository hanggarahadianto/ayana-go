# Gunakan golang sebagai builder
FROM golang:1.21 AS builder

WORKDIR /app

# Hanya install git dan ca-certificates, tidak perlu golang lagi
RUN apt update && apt install -y \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy go.mod dan go.sum sebelum copy semua file
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi dengan path absolut agar tidak ada kesalahan path
RUN GOARCH=amd64 GOOS=linux go build -o /app/main .

# Gunakan Ubuntu sebagai base image
FROM ubuntu:24.04

WORKDIR /app

# Install CA certificates
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

# Copy binary dengan path absolut untuk memastikan file ada
COPY --from=builder /app/main /app/main

# Beri izin eksekusi ke binary
RUN chmod +x /app/main

# Expose port aplikasi
EXPOSE 5000

# Jalankan aplikasi
CMD ["/app/main"]
