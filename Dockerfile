FROM golang:1.21-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main . && chmod +x main

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["/app/main"]
