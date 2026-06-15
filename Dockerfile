# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/main ./cmd

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /src/internal/config/config.default.yaml /app/internal/config/config.default.yaml

EXPOSE 8000

CMD ["/app/main"]
