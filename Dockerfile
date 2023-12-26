# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /docker-gs-ping

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /docker-gs-ping /app/docker-gs-ping

EXPOSE 8080

CMD ["/app/docker-gs-ping"]