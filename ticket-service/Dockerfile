FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.work go.work.sum ./
COPY pkg/ pkg/
COPY proto/ proto/
COPY ticket-service/ ticket-service/

WORKDIR /app/ticket-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /ticket-service ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /ticket-service .
EXPOSE 50053
CMD ["./ticket-service"]
