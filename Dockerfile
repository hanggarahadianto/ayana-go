# Stage 1: Build aplikasi dalam builder container
FROM ubuntu:24.04 AS builder

WORKDIR /app

# Install dependensi yang diperlukan
RUN apt update && apt install -y \
    golang \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Setel environment Go
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Salin file go.mod dan go.sum, lalu download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode aplikasi
COPY . ./

# Bangun binary Go di dalam container
RUN go build -o /app/main .

# Stage 2: Jalankan aplikasi di container final
FROM ubuntu:24.04

WORKDIR /app

# Install CA certificates
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

# Salin binary dari builder stage
COPY --from=builder /app/main /app/main

# Verifikasi file binary ada
RUN ls -alh /app

# Pastikan binary bisa dieksekusi
RUN chmod +x /app/main

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["/app/main"]