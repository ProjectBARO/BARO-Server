FROM golang:1.21-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project (including external packages)
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /docker-gs-ping

EXPOSE 8080

CMD ["/docker-gs-ping"]