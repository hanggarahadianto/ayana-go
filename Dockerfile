# Build stage
FROM golang:1.14.2 as build

WORKDIR /go/src/ayana
COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -v -o ayana

# Run stage
FROM alpine:3.11
COPY --from=build go/src/app/ app/
CMD ["./app/ayana"]