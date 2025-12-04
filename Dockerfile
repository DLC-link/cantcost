# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /log-catcher ./bin/main.go

# Stage 2: Run
FROM alpine:latest

# You can add ca-certificates if your app needs TLS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /log-catcher .

EXPOSE 8080

CMD ["./log-catcher"]