FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.work go.work.sum ./
COPY pkg/ pkg/
COPY proto/ proto/
COPY user-service/ user-service/

WORKDIR /app/user-service
ENV GOWORK=off
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /user-service .
EXPOSE 50051
CMD ["./user-service"]
