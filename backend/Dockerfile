# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build -o sentrinet ./server/main.go

# --- Final lightweight image ---
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/sentrinet /app/sentrinet
COPY sentrinet.db /app/sentrinet.db

EXPOSE 8080
ENTRYPOINT ["/app/sentrinet"]
