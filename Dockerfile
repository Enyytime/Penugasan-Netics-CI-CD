# --- Build Stage ---
FROM golang:1.23-alpine AS builder
WORKDIR /app
# Copy the go.mod and go.sum files and download dependencies
COPY go.mod ./
RUN go mod download
# Copy the source code
COPY . .
# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o health-api .

# --- Final Stage ---
FROM alpine:latest
WORKDIR /root/
# Copy the built binary from the builder stage
COPY --from=builder /app/health-api .
EXPOSE 8080
ENTRYPOINT ["./health-api"]
