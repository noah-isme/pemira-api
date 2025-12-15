# Builder stage
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/api ./cmd/api

# Runner stage
FROM alpine:latest

WORKDIR /

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 1000 nonroot && \
    adduser -D -u 1000 -G nonroot nonroot

COPY --from=builder /bin/api /api
COPY --from=builder /app/migrations /migrations

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/api"]
