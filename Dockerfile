FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-cli ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/chat-cli .
ENTRYPOINT ["./chat-cli"] 