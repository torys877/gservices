FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /validator-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /validator-service .

EXPOSE 8080

CMD ["./validator-service"]
