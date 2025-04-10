# --- Build Stage ---
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o health-api .

# --- Final Stage ---
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/health-api .
EXPOSE 8080
ENTRYPOINT ["./health-api"]
