# -- Build --
FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o lighthouse ./cmd/lighthouse

# -- Final --
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache docker-cli docker-cli-compose

RUN mkdir -p /app/lighthouse

COPY --from=builder /app/lighthouse /lighthouse

CMD ["/lighthouse"]