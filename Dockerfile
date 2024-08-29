# Build stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOARCH=amd64 GOOS=linux go build -o ./bin/img_upload ./cmd/server

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/img_upload .

EXPOSE 8080

CMD ["./img_upload"]
