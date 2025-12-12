# Build Stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD ["/app/main"]
